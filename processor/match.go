package processor

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/gragas/riotclient"
)

// MatchID identifies a match.
type MatchID struct {
	Region string
	ID     int
}

// Matches is the processor for matches.
type Matches struct {
	Riot      *riotclient.RiotClient `inject:"t"`
	Logger    logrus.Logger          `inject:"t"`
	Summoners *Summoners             `inject:"t"`
	c         chan MatchID
}

// NewMatches creates a new processor.Matches.
func NewMatches() *Matches {
	return &Matches{
		c: make(chan MatchID),
	}
}

// Offer offers a match to the queue which may accept it.
func (m *Matches) Offer(id MatchID) {
	// if key exists in cassandra return
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

func (m *Matches) process(id MatchID) {
	region := s.Riot.Region(id.Region)
	res, err := region.Match(strconv.Itoa(m.ID))
	if err != nil {
		s.Logger.Errorf("Could not fetch details of matach %s in region %s: %v", id.ID, id.Region, err)
		return
	}
	// TODO(simplyianm): actually do shit
	for _, p := range res.ParticipantIdentities {
		m.Summoners.Offer(SummonerID{
			Region: id.Region,
			ID:     p.Player.SummonerID,
		})
	}
}
