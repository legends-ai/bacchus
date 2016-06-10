package spider

import (
	"fmt"
	"strconv"

	"github.com/simplyianm/gragas/structures"
	"github.com/simplyianm/gragas/util"
	"github.com/simplyianm/riotclient"
)

const (
	nameChunkSize = 20
)

// Spider represents a spider to search the riot API
type Spider struct {
	Riot        *riotclient.API
	Games       *structures.Queue
	Summoners   *structures.Queue
	Concurrency int
}

func Create(api *riotclient.API, concurrency int) *Spider {
	return &Spider{
		Riot:        api,
		Games:       structures.NewQueue(),
		Summoners:   structures.NewQueue(),
		Concurrency: concurrency,
	}
}

// SeedSummoners seeds the summoners
func (s *Spider) SeedSummoners(summonerIds []string) {
	for _, summoner := range summonerIds {
		s.Summoners.Offer(summoner)
	}
}

// SeedFromFeaturedGames seeds the spider with featured games summoners
func (s *Spider) SeedFromFeaturedGames() error {
	r, err := s.Riot.FeaturedGames()
	if err != nil {
		return fmt.Errorf("Could not get featured games: %v", err)
	}
	names := structures.NewStringSet()
	for _, g := range r.GameList {
		for _, p := range g.Participants {
			names.Add(p.SummonerName)
		}
	}
	return s.seedSummonersByName(names.Values())
}

func (s *Spider) seedSummonersByName(summoners []string) error {
	chunks := util.Chunk(summoners, nameChunkSize)
	for _, chunk := range chunks {
		sum, err := s.Riot.SummonerByName(chunk)
		if err != nil {
			return err
		}
		for _, summoner := range sum {
			s.Summoners.Offer(strconv.Itoa(summoner.ID))
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

func (s *Spider) processGame(g string) {
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
