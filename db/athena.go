// Package db interacts with the database
package db

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/config"
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

// MatchID identifies a match.
type MatchID struct {
	Region string
	ID     int
}

// String returns a string representation of this ID.
func (id MatchID) String() string {
	return fmt.Sprintf("%s/%s", id.Region, id.ID)
}

// Rank represents a rank.
type Rank struct {
	Division uint32
	Tier     uint32
}

// ToNumber returns a numerical representation of rank that can be sorted.
func (r *Rank) ToNumber() uint64 {
	return uint64(r.Tier)<<32 | uint64(r.Division)
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

// WriteMatch writes a match to Cassandra.
func (a *Athena) WriteMatch(m *Match) error {
	return a.Session.Query(
		insertMatchQuery, m.ID.String(),
		m.ID.ID, m.ID.Region, m.Body, m.Rank.ToNumber(),
	).Exec()
}

// Match represents a match.
type Match struct {
	ID   MatchID
	Body string
	Rank Rank
}
