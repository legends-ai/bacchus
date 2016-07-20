package processor

import (
	"encoding/json"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/models"
	"github.com/simplyianm/bacchus/rank"
	"github.com/simplyianm/bacchus/riot"
)

// Matches is the processor for matches.
type Matches struct {
	Riot      *riot.Client        `inject:"t"`
	Logger    *logrus.Logger      `inject:"t"`
	Summoners *Summoners          `inject:"t"`
	Matches   *db.MatchesDAO      `inject:"t"`
	Ranks     *rank.LookupService `inject:"t"`
	c         chan models.MatchID
	cutoff    models.Rank
}

// NewMatches creates a new processor.Matches.
func NewMatches() *Matches {
	cutoff, _ := models.ParseRank(models.TierPlatinum, models.DivisionV)
	return &Matches{
		c:      make(chan models.MatchID, 1E7),
		cutoff: *cutoff,
	}
}

// Offer offers a match to the queue which may accept it.
func (m *Matches) Offer(id models.MatchID) {
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

func (m *Matches) process(id models.MatchID) {
	m.Logger.Infof("Processing match %s", id.String())
	region := m.Riot.Region(id.Region)

	// Retrieve match data
	res, err := region.Match(strconv.Itoa(id.ID))
	m.Logger.Infof("Fetched match data for %s", id.String())
	if err != nil {
		m.Logger.Errorf("Could not fetch details of match %s in region %s: %v", id.ID, id.Region, err)
		return
	}

	// Ignore non-ranked
	m.Logger.Infof("Checking correct queue for %s", id.String())
	if res.QueueType != riot.QueueSolo5x5 && res.QueueType != riot.QueuePremade5x5 && res.QueueType != riot.QueueTeamBuilderDraftRanked5x5 {
		m.Logger.Infof("Wrong queue for %s: %s", id.String(), res.QueueType)
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
	m.Logger.Infof("Getting min ranks for %s", id.String())
	sums := m.Ranks.Lookup(ids, res.Time())

	var ranks []models.Rank
	for id, r := range sums {
		ranks = append(ranks, r)
		if r.Over(m.cutoff) {
			m.Summoners.Offer(id)
		}
	}
	rank := models.MinRank(ranks)

	// Minify JSON
	json, err := m.minifyJSON(res.RawJSON)
	if err != nil {
		m.Logger.Errorf("Could not minify Riot JSON: %v", err)
	}

	m.Logger.Infof("Wrote match %s with rank %d", id.String(), rank.ToNumber())

	// Write match to Cassandra
	m.Matches.Insert(&models.Match{
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
