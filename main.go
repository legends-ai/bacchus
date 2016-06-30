package main

import (
	"github.com/simplyianm/gragas/processor"
	"github.com/simplyianm/gragas/riotclient"
	"github.com/simplyianm/inject"
	"github.com/simplyianm/keypool"
	"github.com/simplyianm/riot/config"
)

const (
	concurrency = 10
)

func main() {
	injector := inject.New()
	injector.Map(injector)

	// Load config
	cfg := config.Fetch()
	injector.Map(cfg)

	// Load keypool
	keys := keypool.New(cfg.APIKeys, cfg.MaxRate)
	injector.Map(keys)

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
