package spider

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/simplyianm/gragas/structures"
	"github.com/simplyianm/gragas/util"
	"github.com/simplyianm/riotclient"
)

const (
	nameChunkSize = 20
)

// Spider represents a spider to search the riot API
type Spider struct {
	Games        *structures.Queue
	Summoners    *structures.Queue
	Concurrency  int
	clients      map[string]*riotclient.API
	clientsMutex sync.Mutex
}

func Create(api *riotclient.API, concurrency int) *Spider {
	return &Spider{
		Games:       structures.NewQueue(),
		Summoners:   structures.NewQueue(),
		Concurrency: concurrency,
		clients:     make(map[string]*riotclient.API),
	}
}

// SeedSummoners seeds the summoners
func (s *Spider) SeedSummoners(summonerIds []string, region string) {
	for _, summoner := range summonerIds {
		s.Summoners.Offer(structures.RegionedString{summoner, region})
	}
}

// SeedFromFeaturedGames seeds the spider with featured games summoners
func (s *Spider) SeedFromFeaturedGames(region string) error {
	r, err := s.riot(region).FeaturedGames()
	if err != nil {
		return fmt.Errorf("Could not get featured games: %v", err)
	}
	names := structures.NewSet()
	for _, g := range r.GameList {
		for _, p := range g.Participants {
			names.Add(p.SummonerName)
		}
	}
	namesStrs := []string{}
	for _, name := range names.Values() {
		namesStrs = append(namesStrs, name.(string))
	}
	return s.seedSummonersByName(namesStrs, region)
}

func (s *Spider) seedSummonersByName(summoners []string, region string) error {
	chunks := util.Chunk(summoners, nameChunkSize)
	for _, chunk := range chunks {
		sum, err := s.riot(region).SummonerByName(chunk)
		if err != nil {
			return err
		}
		for _, summoner := range sum {
			s.Summoners.Offer(structures.RegionedString{strconv.Itoa(summoner.ID), region})
		}
	}
	return nil
}

// Start the spider
func (s *Spider) Start() {
	for i := 0; i < s.Concurrency; i++ {
		go s.process()
	}
}

func (s *Spider) process() {
	for {
		select {
		case g, more := <-s.Games.Channel:
			if !more {
				continue
			}
			s.processGame(g)

		case summoner, more := <-s.Summoners.Channel:
			if !more {
				continue
			}
			s.processSummoner(summoner)
		}
	}
}

func (s *Spider) processGame(g RegionedString) {
	s.Games.Start(g)
	defer s.Games.Complete(g)
	resp, err := s.Riot.Match(g)
	if err != nil {
		// TODO(simplyianm): retry bad games
		return
	}
	// TODO(simplyianm): call league service
	json := resp.RawJSON
	// TODO(simplyianm): store json
	fmt.Printf("scraped %s\n", g)
	fmt.Println(json)
}

func (s *Spider) processSummoner(summoner string) {
	s.Summoners.Start(summoner)
	defer s.Summoners.Complete(summoner)
	resp, err := s.Riot.Game(summoner)
	if err != nil {
		// TODO(simplyianm): retry bad summoners
		return
	}
	for _, g := range resp.Games {
		s.Games.Offer(strconv.Itoa(g.GameID))
	}
}

// riot retrieves the riot api client for the given region
func (s *Spider) riot(region string) *riotclient.API {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	if client := s.clients[region]; client != nil {
		return client
	}
	client := riotclient.New(region)
	s.clients[region] = client
	return client
}
