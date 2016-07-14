package lib

import (
	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/config"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/processor"
	"github.com/simplyianm/bacchus/rank"
	"github.com/simplyianm/bacchus/riotclient"
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
	_, err := injector.ApplyMap(riotclient.New())
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

	return injector
}
