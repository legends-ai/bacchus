// Package db interacts with the database
package db

import (
	"github.com/asunaio/bacchus/config"
	"github.com/gocql/gocql"
)

// NewSession creates a new database session.
func NewSession(cfg *config.AppConfig) (*gocql.Session, error) {
	c := gocql.NewCluster(cfg.CassandraHosts...)
	c.Keyspace = cfg.CassandraKeyspace
	c.ProtoVersion = 3
	return c.CreateSession()
}
