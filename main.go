package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/config"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/processor"
	"github.com/simplyianm/bacchus/riotclient"
	"github.com/simplyianm/inject"
	"github.com/simplyianm/keypool"
)

const (
	concurrency = 10
)

func main() {
	injector := inject.New()
	injector.Map(injector)

	// Load logger
	logger := logrus.New()
	injector.Map(logger)

	// Load config
	cfg := config.Fetch()
	injector.Map(cfg)

	// Load keypool
	keys := keypool.New(cfg.APIKeys, cfg.MaxRate)
	injector.Map(keys)

	// Load Cassandra cluster
	athena, err := db.NewAthena(cfg)
	if err != nil {
		logger.Fatalf("Could not load Athena cluster: %v", err)
	}
	injector.Map(athena)

	// Create a client for Riot
	injector.ApplyMap(riotclient.New())

	// Load summoner and match processors
	s := processor.NewSummoners()
	injector.Map(s)
	m := processor.NewMatches()
	injector.ApplyMap(m)
	injector.Apply(s)

	// Start processing queues
	for i := 0; i < concurrency; i++ {
		go s.Start()
	}
	for i := 0; i < concurrency; i++ {
		go m.Start()
	}

	select {}
}
