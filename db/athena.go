// Package db interacts with the database
package db

import (
	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/config"
	"github.com/simplyianm/bacchus/models"
)

const (
	hasMatchQuery    = `SELECT COUNT(*) FROM matches WHERE id = ?`
	insertMatchQuery = `INSERT INTO matches (id, match_id, region, body, rank) VALUES (?, ?, ?, ?, ?)`
)

// Athena is the athena cluster
type Athena struct {
	Session *gocql.Session
}

// NewAthena creates a new Athena object from config
func NewAthena(cfg *config.AppConfig) (*Athena, error) {
	c := gocql.NewCluster(cfg.AthenaHosts...)
	c.Keyspace = cfg.AthenaKeyspace
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
		insertMatchQuery, m.ID.String(),
		m.ID.ID, m.ID.Region, m.Body, m.Rank.ToNumber(),
	).Exec()
}