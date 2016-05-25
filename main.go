package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/simplyianm/gragas/riot"
	"github.com/simplyianm/gragas/spider"
)

const (
	envRiotApiKey = "RIOT_API_KEY"
	concurrency   = 10
)

func main() {
	apiKey := os.Getenv(envRiotApiKey)
	if apiKey == "" {
		log.Fatalf("Missing %s variable", envRiotApiKey)
	}
	api := riot.APISettings{
		APIKey: apiKey,
		Region: "na",
	}.Create()
	s, err := spider.Create(api, concurrency)
	if err != nil {
		log.Fatalf("Cannot initialize spider: %v", err)
		return
	}
	s.Games.Print()
	s.Summoners.Print()
	r, err := api.Match(s.Games.Unvisited.Values()[0])
	spew.Dump(r)
}
