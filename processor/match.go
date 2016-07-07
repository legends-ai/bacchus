package processor

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/riotclient"
)

// Matches is the processor for matches.
type Matches struct {
	Riot      *riotclient.RiotClient `inject:"t"`
	Logger    logrus.Logger          `inject:"t"`
	Summoners *Summoners             `inject:"t"`
	Athena    *db.Athena             `inject:"t"`
	c         chan db.MatchID
}

// NewMatches creates a new processor.Matches.
func NewMatches() *Matches {
	return &Matches{
		c: make(chan db.MatchID),
	}
}

// Offer offers a match to the queue which may accept it.
func (m *Matches) Offer(id db.MatchID) {
	// if key exists in cassandra return
	ok, err := m.Athena.HasMatch(id)
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

func (m *Matches) process(id db.MatchID) {
	region := m.Riot.Region(id.Region)
	res, err := region.Match(strconv.Itoa(id.ID))
	if err != nil {
		m.Logger.Errorf("Could not fetch details of matach %s in region %s: %v", id.ID, id.Region, err)
		return
	}
	// TODO(simplyianm): actually do shit
	var ids []SummonerID
	for _, p := range res.ParticipantIdentities {
		ids = append(ids, SummonerID{
			Region: id.Region,
			ID:     p.Player.SummonerID,
		})
	}
}
