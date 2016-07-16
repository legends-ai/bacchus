package config

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	app = "RIOT"
)

var uuidMatcher = regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$")

// AppConfig is the configuration for the app
type AppConfig struct {
	// APIKeys is a list of API keys to pool
	APIKeys []string `envconfig:"api_keys" required:"true"`

	// MaxRate is the maximum rate to make requests per key
	// The default corresponds to 500 requests per 10 minutes.
	MaxRate time.Duration `envconfig:"max_rate" default:"1200ms"`

	// AthenaHosts is a list of Athena Cassandra hosts.
	AthenaHosts []string `envconfig:"athena_hosts" required:"true"`

	// RankExpiry is the max duration since a rank is valid.
	RankExpiry time.Duration `envconfig:"rank_expiry" default:"168h"`

	// BatchSize is the size of a batch when performing rank lookups.
	BatchSize int `envconfig:"batch_size" default:"8"`
}

// Fetch fetches the config from env vars
func Fetch() *AppConfig {
	var a AppConfig
	err := envconfig.Process(app, &a)
	if err != nil {
		log.Fatalf("Error processing config: %v", err)
	}
	err = a.Validate()
	if err != nil {
		log.Fatalf("Error validating config: %v", err)
	}
	return &a
}

// Validate validates the app config
func (a *AppConfig) Validate() error {
	for _, key := range a.APIKeys {
		if !isValidKey(key) {
			return fmt.Errorf("Invalid api key: %s", key)
		}
	}
	return nil
}

func isValidKey(key string) bool {
	return uuidMatcher.MatchString(key) || strings.HasPrefix(key, "RGAPI")
}
