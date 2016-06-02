package main

import (
	"log"

	"github.com/simplyianm/gragas/riot"
	"github.com/simplyianm/gragas/spider"
)

const (
	envRiotApiKey = "RIOT_API_KEY"
	concurrency   = 10
)

func main() {
	api := riot.APISettings{
		Region: "na",
	}.Create()
	s, err := spider.Create(api, concurrency)
	if err != nil {
		log.Fatalf("Cannot initialize spider: %v", err)
		return
	}
	s.Start()
	select {}
}
