package main

import (
	"log"

	"github.com/simplyianm/gragas/spider"
	"github.com/simplyianm/riotclient"
)

const (
	concurrency = 10
)

func main() {
	api := riotclient.Create("na")
	s := spider.Create(api, 10)
	if err := s.SeedFromFeaturedGames(); err != nil {
		log.Fatalf("Cannot seed spider: %v", err)
		return
	}
	s.Start()
	select {}
}
