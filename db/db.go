// Package db interacts with the database
package db

import (
	"github.com/gocql/gocql"
	"github.com/simplyianm/bacchus/config"
)

const keyspace = "athena"

// NewSession creates a new database session.
func NewSession(cfg *config.AppConfig) (*gocql.Session, error) {
	c := gocql.NewCluster(cfg.AthenaHosts...)
	c.Keyspace = keyspace
	c.ProtoVersion = 3
	return c.CreateSession()
}
