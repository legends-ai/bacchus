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
	s, err := spider.Create(api, concurrency)
	if err != nil {
		log.Fatalf("Cannot initialize spider: %v", err)
		return
	}
	s.Start()
	select {}
}
