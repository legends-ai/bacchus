package main

import (
	"log"
	"os"

	"github.com/simplyianm/gragas/clients"
	"github.com/simplyianm/gragas/spider"
)

func main() {
	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		log.Fatalf("Missing RIOT_API_KEY variable")
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
