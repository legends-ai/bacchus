package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	r, _ := api.Game(s.Summoners.Unvisited.Values()[0])
	games := r.Games
	gameIds := []int{}
	for _, g := range games {
		gameIds = append(gameIds, g.GameId)
	}
	x, _ := api.Match(strconv.Itoa(gameIds[0]))
	fmt.Println(x.RawJSON)
}
