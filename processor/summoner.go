package processor

import (
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/models"
	"github.com/simplyianm/bacchus/riot"
)

// Queues is the processor for queues.
type Summoners struct {
	Riot     *riot.Client    `inject:"t"`
	Logger   *logrus.Logger  `inject:"t"`
	Matches  *Matches        `inject:"t"`
	Rankings *db.RankingsDAO `inject:"t"`
	c        chan models.SummonerID
	exists   map[models.SummonerID]bool
	existsMu sync.RWMutex
}

// NewSummoners creates a new processor.Summoners.
func NewSummoners() *Summoners {
	return &Summoners{
		c:      make(chan models.SummonerID, 1E7),
		exists: map[models.SummonerID]bool{},
	}
}

// Offer offers a summoner to the queue which may accept it.
func (s *Summoners) Offer(id models.SummonerID) {
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
	rank, err := models.ParseRank(models.TierPlatinum, models.DivisionV)
	if err != nil {
		s.Logger.Fatalf("Invalid static rank: %v", err)
	}
	ids, err := s.Rankings.AboveRank(*rank, 1000)
	if err != nil {
		s.Logger.Fatalf("Could not perform initial seed: %v", err)
	}
	if len(ids) == 0 {
		// no plat, do alternative seed
		s.Offer(models.SummonerID{"na", 32875076}) // Pradyuman himself
		return
	}
	for _, id := range ids {
		s.Offer(id)
	}
}

func (s *Summoners) process(id models.SummonerID) {
	s.Logger.Infof("Processing summoner %s", id.String())
	res, err := s.Riot.Region(id.Region).Game(strconv.Itoa(id.ID))
	if err != nil {
		s.Logger.Errorf("Could not fetch games of summoner %s in region %s: %v", id.ID, id.Region, err)
		return
	}
	s.existsMu.Lock()
	s.exists[id] = true
	s.existsMu.Unlock()
	for _, game := range res.Games {
		s.Matches.Offer(models.MatchID{
			Region: id.Region,
			ID:     game.GameID,
		})
	}
}
