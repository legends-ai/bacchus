package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/config"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/models"
	"github.com/simplyianm/bacchus/processor"
	"github.com/simplyianm/bacchus/rank"
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

	// Create a client for Riot
	_, err := injector.ApplyMap(riotclient.New())
	if err != nil {
		logger.Fatalf("Could not inject riot client: %v", err)
	}

	// Load Cassandra cluster
	logger.Info("Connecting to Athena Cassandra")
	athena, err := db.NewAthena(cfg)
	if err != nil {
		logger.Fatalf("Could not load Athena cluster: %v", err)
	}
	injector.Map(athena)

	// Load lookup service
	_, err = injector.ApplyMap(&rank.LookupService{})
	if err != nil {
		logger.Fatalf("Could not inject lookup service: %v", err)
	}

	// Load summoner and match processors
	logger.Info("Loading processors")
	s := processor.NewSummoners()
	injector.Map(s)
	m := processor.NewMatches()

	injector.ApplyMap(m)
	if err != nil {
		logger.Fatalf("Could not inject match processor: %v", err)
	}

	injector.Apply(s)
	if err != nil {
		logger.Fatalf("Could not inject summoner processor: %v", err)
	}

	// Start processing queues
	for i := 0; i < concurrency; i++ {
		logger.Infof("Starting summoner processor %d", i)
		go s.Start()
	}
	for i := 0; i < concurrency; i++ {
		logger.Infof("Starting match processor %d", i)
		go m.Start()
	}

	// Offer Aditi
	s.Offer(models.SummonerID{"na", 32875076})

	select {}
}
