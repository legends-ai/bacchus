// Package db interacts with the database
package db

import (
	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/config"
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
