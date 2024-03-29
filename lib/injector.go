package lib

import (
	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/inject"
	"google.golang.org/grpc"
	"gopkg.in/redis.v5"

	"github.com/asunaio/bacchus/config"
	"github.com/asunaio/bacchus/db"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/processor"
	"github.com/asunaio/bacchus/queue"
	"github.com/asunaio/bacchus/rank"
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

	// Charon Client
	logger.Infof("Connecting to Charon at %s", cfg.CharonHost)
	charonConn, err := grpc.Dial(cfg.CharonHost, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("Could not connect to Charon: %v", err)
	}
	charon := apb.NewCharonClient(charonConn)
	injector.Map(charon)

	// Totsuki Client
	logger.Infof("Connecting to Totsuki at %s", cfg.TotsukiHost)
	totsukiConn, err := grpc.Dial(cfg.TotsukiHost, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("Could not connect to Totsuki: %v", err)
	}
	totsuki := apb.NewTotsukiClient(totsukiConn)
	injector.Map(totsuki)

	// Load Cassandra cluster
	logger.Info("Connecting to Cassandra")
	session, err := db.NewSession(cfg)
	if err != nil {
		logger.Fatalf("Could not load Cassandra cluster: %v", err)
	}
	injector.Map(session)

	// Redis Client
	logger.Infof("Connecting to Redis at %s", cfg.RedisHost)
	redisConn := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: "",
		DB:       0,
		PoolSize: 10,
	})
	_, err = redisConn.Ping().Result()
	if err != nil {
		logger.Fatalf("Could not connect to Redis: %v", err)
	}
	injector.Map(redisConn)

	// Initialize queues
	_, err = injector.ApplyMap(queue.NewMatchQueue())
	if err != nil {
		logger.Fatalf("Could not inject match queue: %v", err)
	}
	_, err = injector.ApplyMap(queue.NewSummonerQueue())
	if err != nil {
		logger.Fatalf("Could not inject summoner queue: %v", err)
	}

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
	metrics := &processor.Metrics{}
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
