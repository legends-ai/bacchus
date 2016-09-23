package processor

import (
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/asunaio/bacchus/db"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
	"github.com/asunaio/bacchus/riot"
)

// Queues is the processor for queues.
type Summoners struct {
	Riot     *riot.Client    `inject:"t"`
	Logger   *logrus.Logger  `inject:"t"`
	Matches  *Matches        `inject:"t"`
	Rankings *db.RankingsDAO `inject:"t"`
	Metrics  *Metrics        `inject:"t"`

	c        chan *apb.SummonerId
	exists   map[*apb.SummonerId]bool
	existsMu sync.RWMutex
}

// NewSummoners creates a new processor.Summoners.
func NewSummoners() *Summoners {
	return &Summoners{
		c:      make(chan *apb.SummonerId, 1E7),
		exists: map[*apb.SummonerId]bool{},
	}
}

// Offer offers a summoner to the queue which may accept it.
func (s *Summoners) Offer(id *apb.SummonerId) {
	s.existsMu.RLock()
	defer s.existsMu.RUnlock()
	if s.exists[id] {
		return
	}
	s.c <- id
}

// Start starts processing summoners.
func (s *Summoners) Start() {
	for {
		id, ok := <-s.c
		if !ok {
			return
		}
		s.process(id)
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
	ids, err := s.Rankings.AboveRank(rank, 1000)
	if err != nil {
		s.Logger.Fatalf("Could not perform initial seed: %v", err)
	}

	// Seed ids
	if len(ids) == 0 {
		// no plat, do alternative seed
		s.Offer(&apb.SummonerId{
			Region: apb.Region_NA,
			Id:     32875076,
		}) // Pradyuman himself
		return
	}

	// Offer found ids
	for _, id := range ids {
		s.Offer(id)
	}
}

func (s *Summoners) process(id *apb.SummonerId) {
	s.Logger.Infof("Processing summoner %s", id.String())

	// process the summoner
	res, err := s.Riot.Region(id.Region).Game(strconv.Itoa(int(id.Id)))
	if err != nil {
		s.Logger.Errorf("Could not fetch games of %s: %v", id.String(), err)
		return
	}

	// lock
	s.existsMu.Lock()
	s.exists[id] = true
	s.existsMu.Unlock()

	// offer le games
	for _, game := range res.Games {
		s.Matches.Offer(&apb.MatchId{
			Region: id.Region,
			Id:     uint32(game.GameID),
		})
	}
}
