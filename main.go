package main

import (
	"log"

	"github.com/simplyianm/gragas/riot"
	"github.com/simplyianm/gragas/spider"
)

const (
	concurrency = 10
)

func main() {
	api := riot.Create("na")
	s, err := spider.Create(api, concurrency)
	if err != nil {
		log.Fatalf("Cannot initialize spider: %v", err)
		return
	}
	s.Start()
	select {}
}
