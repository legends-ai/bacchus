// Package db interacts with the database
package db

import (
	"github.com/asunaio/bacchus/config"
	"github.com/gocql/gocql"
)

const keyspace = "athena"

// NewSession creates a new database session.
func NewSession(cfg *config.AppConfig) (*gocql.Session, error) {
	c := gocql.NewCluster(cfg.AthenaHosts...)
	c.Keyspace = keyspace
	c.ProtoVersion = 3
	return c.CreateSession()
}
