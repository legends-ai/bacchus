// Package db interacts with the database
package db

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/config"
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

// MatchID identifies a match.
type MatchID struct {
	Region string
	ID     int
}

// String returns a string representation of this ID.
func (id MatchID) String() string {
	return fmt.Sprintf("%s/%s", id.Region, id.ID)
}

// SummonerID identifies a summoner.
type SummonerID struct {
	Region string
	ID     int
}

// String returns a string representation of this ID.
func (id SummonerID) String() string {
	return fmt.Sprintf("%s/%s", id.Region, id.ID)
}

// HasMatch returns true if a match exists.
func (a *Athena) HasMatch(id MatchID) (bool, error) {
	var count int
	if err := a.Session.Query(hasMatchQuery, id.String()).Scan(&count); err != nil {
		return false, err
	}
	return count != 0, nil
}
