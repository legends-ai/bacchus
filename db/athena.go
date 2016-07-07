// Package db interacts with the database
package db

import (
	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/config"
	"github.com/simplyianm/bacchus/processor"
)

const (
	hasMatchQuery = `SELECT COUNT(*) FROM matches WHERE id = ?`
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
func (a *Athena) HasMatch(id processor.MatchID) (bool, error) {
	var count int
	if err := a.Session.Query(hasMatchQuery, id.String()).Scan(&count); err != nil {
		return nil, err
	}
	return count != 0, nil
}
