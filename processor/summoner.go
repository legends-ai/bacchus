package processor

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/gragas/riotclient"
)

// SummonerID identifies a summoner.
type SummonerID struct {
	Region string
	ID     int
}

// Queues is the processor for queues.
type Summoners struct {
	Riot    *riotclient.RiotClient `inject:"t"`
	Logger  logrus.Logger          `inject:"t"`
	Matches *Matches               `inject:"t"`
	c       chan SummonerID
	exists  map[SummonerID]bool
}

// NewSummoners creates a new processor.Summoners.
func NewSummoners() *Summoners {
	return &Summoners{
		c:      make(chan SummonerID),
		exists: map[SummonerID]bool{},
	}
}

// Offer offers a summoner to the queue which may accept it.
func (s *Summoners) Offer(id SummonerID) {
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

func (s *Summoners) process(m SummonerID) {
	res, err := s.Riot.Region(m.Region).Game(strconv.Itoa(m.ID))
	if err != nil {
		s.Logger.Errorf("Could not fetch games of summoner %s in region %s: %v", m.ID, m.Region, err)
		return
	}
	for _, game := range res.Games {
		s.Matches.Offer(MatchID{
			Region: m.Region,
			ID:     game.GameID,
		})
	}
}
