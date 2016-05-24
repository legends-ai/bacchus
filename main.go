package main

import (
	"log"
	"os"

	"github.com/simplyianm/gragas/clients"
	"github.com/simplyianm/gragas/spider"
)

const (
	envRiotApiKey = "RIOT_API_KEY"
)

func main() {
	apiKey := os.Getenv(envRiotApiKey)
	if apiKey == "" {
		log.Fatalf("Missing %s variable", apiKey)
	}
	api := clients.RiotAPISettings{
		APIKey: apiKey,
		Region: "na",
	}.Create()
	s, err := spider.Create(api)
	if err != nil {
		log.Fatalf("Cannot initialize spider: %v", err)
		return
	}
	s.Games.Print()
}
