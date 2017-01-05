package processor

import (
	"github.com/Shopify/sarama"
	"github.com/Sirupsen/logrus"
	"github.com/asunaio/bacchus/db"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
	"github.com/asunaio/bacchus/queue"
	"github.com/asunaio/bacchus/rank"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// Matches is the processor for matches.
type Matches struct {
	Charon    apb.CharonClient     `inject:"t"`
	Logger    *logrus.Logger       `inject:"t"`
	Matches   *db.MatchesDAO       `inject:"t"`
	Kafka     sarama.AsyncProducer `inject:"t"`
	Metrics   *Metrics             `inject:"t"`
	Ranks     *rank.LookupService  `inject:"t"`
	Summoners *Summoners           `inject:"t"`
	Queue     *queue.MatchQueue    `inject:"t"`

	cutoff *apb.Rank
}

// NewMatches creates a new processor.Matches.
func NewMatches() *Matches {
	cutoff, _ := models.ParseRank(models.TierPlatinum, models.DivisionV)
	return &Matches{
		cutoff: cutoff,
	}
}

// Offer offers a match to the queue which may accept it.
func (m *Matches) Offer(info *apb.CharonRpc_MatchListResponse_MatchInfo) {
	// if key exists in cassandra return
	if ok, err := m.Matches.Exists(info.MatchId); err != nil {
		m.Logger.Warnf("Could not check match: %v", err)
		return
	} else if ok {
		m.Metrics.Record("match-duplicates")
		return
	}

	m.Queue.Add(info.MatchId, info)
}

// Start starts processing matches.
func (m *Matches) Start() {

	// Poll for producer errors
	go func() {
		for err := range m.Kafka.Errors() {
			m.Logger.Errorf("Matches Producer Error: %v", err)
		}
	}()

	for {
		m.process(m.Queue.Poll())
	}
}

func (m *Matches) process(id *apb.MatchId) {
	// Retrieve match data
	res, err := m.Charon.GetMatch(context.TODO(), &apb.CharonRpc_MatchRequest{
		Match: id,
	})
	if err != nil {
		m.Logger.Errorf("Could not fetch details of match %v: %v", id, err)
		return
	}

	m.Metrics.Record("match-fetch")

	// Fetch summoners from match
	ids := res.Summoners

	// Get min rank of players
	sums, err := m.Ranks.Lookup(ids)
	if err != nil {
		m.Logger.Errorf("Error looking up ranks: %v", err)
		return
	}

	var ranks []*apb.Rank
	for _, ranking := range sums {
		ranks = append(ranks, ranking.Rank)
		if models.RankOver(ranking.Rank, m.cutoff) {
			m.Summoners.Offer(ranking)
		}
	}

	rank := models.MedianRank(ranks)
	if rank == nil {
		m.Logger.Errorf("Outdated ranks for match %s", id)
		return
	}

	match, err := proto.Marshal(&apb.BacchusData_RawMatch{
		Id:      id,
		Version: res.MatchInfo.Version,
		Rank:    rank,
		Data:    res.MatchInfo,
	})

	if err != nil {
		m.Logger.Errorf("Error marshaling match: %v", err)
		return
	}

	// Publish match record to Kafka
	m.Kafka.Input() <- &sarama.ProducerMessage{
		Topic: "bacchus.matches." + id.Region.String(),
		Value: sarama.ByteEncoder(match),
	}

	// Write match id to Cassandra
	if err := m.Matches.Insert(id); err != nil {
		m.Logger.Errorf("Could not insert match id to Cassandra: %v", err)
	}

	m.Metrics.Record("match-write")
}
