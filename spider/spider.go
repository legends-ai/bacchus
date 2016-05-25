package spider

import (
	"fmt"

	"github.com/simplyianm/gragas/riot"
	"github.com/simplyianm/gragas/structures"
)

type Spider struct {
	Riot        *riot.API
	Games       *structures.Queue
	Summoners   *structures.Queue
	Concurrency int
}

func Create(api *riot.API, concurrency int) (*Spider, error) {
	s := &Spider{
		Riot:        api,
		Games:       structures.NewQueue(),
		Summoners:   structures.NewQueue(),
		Concurrency: concurrency,
	}
	err := s.seedFromFeaturedGames()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Spider) seedFromFeaturedGames() error {
	r, err := s.Riot.FeaturedGames()
	if err != nil {
		return fmt.Errorf("Could not get featured games: %v", err)
	}
	for _, g := range r.GameList {
		for _, p := range g.Participants {
			s.Summoners.Offer(p.SummonerName)
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
}

func (s *Spider) processSummoner(summoner string) {
	s.Summoners.Start(summoner)
	defer s.Summoners.Complete(summoner)
}
