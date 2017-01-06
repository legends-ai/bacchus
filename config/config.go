package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// AppConfig is the configuration for the app
type AppConfig struct {
	// CassandraHosts is a list of Cassandra hosts.
	CassandraHosts []string `envconfig:"cassandra_hosts" default:"localhost:9042"`

	// CassandraKeyspace is the Bacchus Cassandra keyspace.
	CassandraKeyspace string `envconfig:"cassandra_keyspace" default:"bacchus_dev"`

	// RankExpiry is the max duration since a rank is valid.
	RankExpiry time.Duration `envconfig:"rank_expiry" default:"168h"`

	// BatchSize is the size of a batch when performing rank lookups.
	BatchSize int `envconfig:"batch_size" default:"8"`

	// Port is the port on which Bacchus runs on.
	Port int `envconfig:"port" default:"9730"`

	// MonitorPort is the port on which pprof and health check run on.
	MonitorPort int `envconfig:"monitor_port" default:"9731"`

	// CharonHost is the Charon host/port.
	CharonHost string `envconfig:"charon_host" default:"localhost:5609"`

	// RedisHost is the Redis host for the queue
	RedisHost string `envconfig:"redis_host" default:"localhost:6379"`

	// TotsukiHost is the Totsuki host/port.
	TotsukiHost string `envconfig:"totsuki_host" default:"localhost:21215"`

	// Concurrency is the number of parallel threads
	Concurrency int `envconfig:"concurrency" default:"60"`
}

// Fetch fetches the config from env vars
func Fetch() *AppConfig {
	var a AppConfig
	err := envconfig.Process("BACCHUS", &a)
	if err != nil {
		log.Fatalf("Error processing config: %v", err)
	}
	return &a
}
