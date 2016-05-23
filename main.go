package main

import (
	"log"
	"os"

	"github.com/simplyianm/gragas/clients"
	"github.com/simplyianm/gragas/spider"
)

func main() {
	api := clients.RiotAPISettings{
		APIKey: os.Getenv("RIOT_API_KEY"),
		Region: "na",
	}.Create()
	s, err := spider.Create(api)
	if err != nil {
		log.Fatalf("Cannot initialize spider: %v", err)
		return
	}
	s.Games.Print()
}
