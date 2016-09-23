package lib

import (
	"github.com/Sirupsen/logrus"
	"github.com/asunaio/bacchus/config"
	"github.com/asunaio/bacchus/db"
	"github.com/asunaio/bacchus/processor"
	"github.com/asunaio/bacchus/rank"
	"github.com/asunaio/bacchus/riot"
	"github.com/simplyianm/inject"
	"github.com/simplyianm/keypool"
)

// NewInjector sets up dependencies for Bacchus.
func NewInjector() inject.Injector {
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
	_, err := injector.ApplyMap(riot.New())
	if err != nil {
		logger.Fatalf("Could not inject riot client: %v", err)
	}

	// Load Cassandra cluster
	logger.Info("Connecting to Cassandra")
	session, err := db.NewSession(cfg)
	if err != nil {
		logger.Fatalf("Could not load Cassandra cluster: %v", err)
	}
	injector.Map(session)

	// DAOs
	injector.ApplyMap(&db.MatchesDAO{})
	injector.ApplyMap(&db.RankingsDAO{})

	// Load batcher
	_, err = injector.ApplyMap(rank.NewBatcher())
	if err != nil {
		logger.Fatalf("Could not inject batcher: %v", err)
	}

	// Load lookup service
	_, err = injector.ApplyMap(&rank.LookupService{})
	if err != nil {
		logger.Fatalf("Could not inject lookup service: %v", err)
	}

	// Load processor metrics
	metrics := processor.Metrics{
		SummonerRate: 1,
		MatchRate:    1,
	}
	_, err = injector.ApplyMap(metrics)
	if err != nil {
		logger.Fatalf("Could not inject processor: %v", err)
	}
	go metrics.Start()

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

	return injector
}
