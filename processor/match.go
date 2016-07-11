package processor

import (
	"encoding/json"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/models"
	"github.com/simplyianm/bacchus/rank"
	"github.com/simplyianm/bacchus/riotclient"
)

// Matches is the processor for matches.
type Matches struct {
	Riot      *riotclient.RiotClient `inject:"t"`
	Logger    logrus.Logger          `inject:"t"`
	Summoners *Summoners             `inject:"t"`
	Athena    *db.Athena             `inject:"t"`
	Ranks     *rank.LookupService    `inject:"t"`
	c         chan models.MatchID
}

// NewMatches creates a new processor.Matches.
func NewMatches() *Matches {
	return &Matches{
		c: make(chan models.MatchID),
	}
}

// Offer offers a match to the queue which may accept it.
func (m *Matches) Offer(id models.MatchID) {
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

func (m *Matches) process(id models.MatchID) {
	m.Logger.Info("Processing match %s", id.String())
	region := m.Riot.Region(id.Region)

	// Retrieve match data
	res, err := region.Match(strconv.Itoa(id.ID))
	if err != nil {
		m.Logger.Errorf("Could not fetch details of match %s in region %s: %v", id.ID, id.Region, err)
		return
	}

	// Ignore non-ranked
	if res.QueueType != riotclient.QueueSolo5x5 && res.QueueType != riotclient.QueuePremade5x5 {
		return
	}

	// Fetch summoners from match
	var ids []models.SummonerID
	for _, p := range res.ParticipantIdentities {
		ids = append(ids, models.SummonerID{
			Region: id.Region,
			ID:     p.Player.SummonerID,
		})
	}

	// Get min rank of players
	rank := m.Ranks.MinRank(ids, res.Time())

	// Minify JSON
	json, err := m.minifyJSON(res.RawJSON)
	if err != nil {
		m.Logger.Errorf("Could not minify Riot JSON: %v", err)
	}

	// Write match to Cassandra
	m.Athena.WriteMatch(&models.Match{
		ID:    id,
		Body:  json,
		Patch: res.MatchVersion,
		Rank:  rank,
	})
}

func (m *Matches) minifyJSON(data string) (string, error) {
	var min interface{}
	if err := json.Unmarshal([]byte(data), &min); err != nil {
		return "", err
	}
	d, err := json.Marshal(min)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
