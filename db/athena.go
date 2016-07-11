// Package db interacts with the database
package db

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/config"
	"github.com/simplyianm/bacchus/models"
)

const (
	keyspace         = "athena"
	hasMatchQuery    = `SELECT COUNT(*) FROM matches WHERE id = ?`
	insertMatchQuery = `INSERT INTO matches (id, region, body, rank, patch) VALUES (?, ?, ?, ?, ?)`
	rankingsQuery    = `SELECT rankings FROM rankings WHERE id = ?`
)

// Athena is the athena cluster
type Athena struct {
	Session *gocql.Session
}

// NewAthena creates a new Athena object from config
func NewAthena(cfg *config.AppConfig) (*Athena, error) {
	c := gocql.NewCluster(cfg.AthenaHosts...)
	c.Keyspace = keyspace
	c.ProtoVersion = 3
	s, err := c.CreateSession()
	if err != nil {
		return nil, err
	}
	return &Athena{s}, nil
}

// HasMatch returns true if a match exists.
func (a *Athena) HasMatch(id models.MatchID) (bool, error) {
	var count int
	if err := a.Session.Query(hasMatchQuery, id.String()).Scan(&count); err != nil {
		return false, err
	}
	return count != 0, nil
}

// WriteMatch writes a match to Cassandra.
func (a *Athena) WriteMatch(m *models.Match) error {
	return a.Session.Query(
		insertMatchQuery, m.ID.String(), m.ID.Region,
		m.Body, m.Rank.ToNumber(), m.Patch,
	).Exec()
}

// Rankings grabs all rankings of a summoner.
func (a *Athena) Rankings(id models.SummonerID) (*models.RankingList, error) {
	var rankings []struct {
		Time time.Time
		Rank int64
	}
	if err := a.Session.Query(rankingsQuery, id.String()).Scan(&rankings); err != nil {
		return nil, err
	}
	var ret []*models.Ranking
	for _, ranking := range rankings {
		ret = append(ret, &models.Ranking{
			Time: ranking.Time,
			Rank: models.RankFromNumber(uint32(ranking.Rank)),
		})
	}
	return models.NewRankingList(ret), nil
}
