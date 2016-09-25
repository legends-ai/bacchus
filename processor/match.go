package processor

import (
	"github.com/Sirupsen/logrus"
	"github.com/asunaio/bacchus/db"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
	"github.com/asunaio/bacchus/rank"
	"golang.org/x/net/context"
)

// Matches is the processor for matches.
type Matches struct {
	Charon    apb.CharonClient    `inject:"t"`
	Logger    *logrus.Logger      `inject:"t"`
	Matches   *db.MatchesDAO      `inject:"t"`
	Metrics   *Metrics            `inject:"t"`
	Ranks     *rank.LookupService `inject:"t"`
	Summoners *Summoners          `inject:"t"`

	c      chan *apb.MatchId
	cutoff *apb.Rank
}

// NewMatches creates a new processor.Matches.
func NewMatches() *Matches {
	cutoff, _ := models.ParseRank(models.TierPlatinum, models.DivisionV)
	return &Matches{
		c:      make(chan *apb.MatchId, 1E7),
		cutoff: cutoff,
	}
}

// Offer offers a match to the queue which may accept it.
func (m *Matches) Offer(id *apb.MatchId) {
	// if key exists in cassandra return
	ok, err := m.Matches.Exists(id)
	if err != nil {
		m.Logger.Warnf("Could not check match: %v", err)
		return
	}
	if ok {
		// don't scrape duplicate matches
		return
	}
	m.c <- id
}

// Start starts processing matches.
func (m *Matches) Start() {
	for {
		id, ok := <-m.c
		if !ok {
			return
		}
		m.process(id)
	}
}

func (m *Matches) process(id *apb.MatchId) {
	// Retrieve match data
	res, err := m.Charon.GetMatch(context.TODO(), &apb.CharonMatchRequest{
		Match: id,
	})
	if err != nil {
		m.Logger.Errorf("Could not fetch details of match %v: %v", id, err)
		return
	}

	// Fetch summoners from match
	ids := res.Payload.Summoners

	// Get min rank of players
	sums, err := m.Ranks.Lookup(ids)
	if err != nil {
		m.Logger.Errorf("Error looking up ranks: %v", err)
		return
	}

	var ranks []*apb.Rank
	for id, r := range sums {
		ranks = append(ranks, r)
		if models.RankOver(r, m.cutoff) {
			m.Summoners.Offer(id)
		}
	}
	rank := models.MinRank(ranks)

	// Write match to Cassandra
	if err := m.Matches.Insert(&apb.RawMatch{
		Id:    id,
		Patch: res.Payload.MatchVersion,
		Rank:  rank,
		Body:  res.Payload.RawJson,
	}); err != nil {
		m.Logger.Errorf("Could not insert match to Cassandra: %v", err)
	}

	m.Metrics.RecordMatch(id)
}
