package processor

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/asunaio/bacchus/db"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
	"github.com/asunaio/bacchus/queue"
)

// Queues is the processor for queues.
type Summoners struct {
	Charon   apb.CharonClient     `inject:"t"`
	Logger   *logrus.Logger       `inject:"t"`
	Matches  *Matches             `inject:"t"`
	Metrics  *Metrics             `inject:"t"`
	Rankings *db.RankingsDAO      `inject:"t"`
	Queue    *queue.SummonerQueue `inject:"t"`
}

// NewSummoners creates a new processor.Summoners.
func NewSummoners() *Summoners {
	return &Summoners{}
}

// Offer offers a summoner to the queue which may accept it.
func (s *Summoners) Offer(ranking *apb.Ranking) {
	s.Queue.Add(ranking.Summoner, ranking)
}

// Start starts processing summoners.
func (s *Summoners) Start() {
	for {
		s.process(s.Queue.Poll())
	}
}

// Seed adds summoners from the database based off of rank.
func (s *Summoners) Seed() {
	// setup rank
	rank, err := models.ParseRank(models.TierPlatinum, models.DivisionV)
	if err != nil {
		s.Logger.Fatalf("Invalid static rank: %v", err)
	}

	// get rankings
	rankings, err := s.Rankings.AboveRank(rank, 1000)
	if err != nil {
		s.Logger.Fatalf("Could not perform initial seed: %v", err)
	}

	// Seed rankings
	if len(rankings) == 0 {
		s.Logger.Info("Database empty; seeding with hardcoded value")
		// no plat, do alternative seed
		s.Offer(&apb.Ranking{
			Summoner: &apb.SummonerId{
				Region: apb.Region_NA,
				Id:     32875076,
			},
			Rank: rank,
		}) // Pradyuman himself
		return
	}

	// Offer found rankings
	for _, ranking := range rankings {
		s.Offer(ranking)
	}
}

func (s *Summoners) process(id *apb.SummonerId) {
	// process the summoner
	res, err := s.Charon.GetMatchList(context.TODO(), &apb.CharonRpc_MatchListRequest{
		Summoner: id,
		Seasons: []string{
			"PRESEASON2015",
			"SEASON2015",
			"PRESEASON2016",
			"SEASON2016",
			"PRESEASON2017",
			"SEASON2017",
		},
	})
	if err != nil {
		s.Logger.Errorf("Could not fetch games of %s: %v", id.String(), err)
		return
	}

	// offer les games
	for _, match := range res.Matches {
		s.Matches.Offer(match)
	}

	s.Metrics.Record("summoner")
}
