package processor

import (
	"encoding/json"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/asunaio/bacchus/db"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
	"github.com/asunaio/bacchus/rank"
	"github.com/asunaio/bacchus/riot"
)

// Matches is the processor for matches.
type Matches struct {
	Riot      *riot.Client        `inject:"t"`
	Logger    *logrus.Logger      `inject:"t"`
	Summoners *Summoners          `inject:"t"`
	Matches   *db.MatchesDAO      `inject:"t"`
	Ranks     *rank.LookupService `inject:"t"`

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

func (m *Matches) process(id *apb.MatchId) {
	m.Logger.Infof("Processing match %v", id)
	region := m.Riot.Region(id.Region)

	// Retrieve match data
	res, err := region.Match(strconv.Itoa(int(id.Id)))
	m.Logger.Infof("Fetched match data for %v", id)
	if err != nil {
		m.Logger.Errorf("Could not fetch details of match %v: %v", id, err)
		return
	}

	// Ignore non-ranked
	m.Logger.Infof("Checking correct queue for %v", id)
	if res.QueueType != riot.QueueSolo5x5 && res.QueueType != riot.QueuePremade5x5 && res.QueueType != riot.QueueTeamBuilderDraftRanked5x5 {
		m.Logger.Infof("Wrong queue for %v: %s", id, res.QueueType)
		return
	}

	// Fetch summoners from match
	var ids []*apb.SummonerId
	for _, p := range res.ParticipantIdentities {
		ids = append(ids, &apb.SummonerId{
			Region: id.Region,
			Id:     uint32(p.Player.SummonerID),
		})
	}

	// Get min rank of players
	m.Logger.Infof("Getting min ranks for %v", id)
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

	// Minify JSON
	json, err := m.minifyJSON(res.RawJSON)
	if err != nil {
		m.Logger.Errorf("Could not minify Riot JSON: %v", err)
	}

	m.Logger.Infof("Wrote match %v with rank %d", id, models.RankToNumber(rank))

	// Write match to Cassandra
	m.Matches.Insert(&apb.RawMatch{
		Id:    id,
		Patch: res.MatchVersion,
		Rank:  rank,
		Body:  json,
	})
}
