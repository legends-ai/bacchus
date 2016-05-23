package spider

import (
	"fmt"
	"strconv"

	"github.com/simplyianm/gragas/clients"
	"github.com/simplyianm/gragas/structures"
)

type Spider struct {
	Riot      *clients.RiotAPI
	Games     *structures.Queue
	Summoners *structures.Queue
}

func Create(api *clients.RiotAPI) (*Spider, error) {
	s := &Spider{
		Riot: api,
		Games: structures.QueueSettings{
			Concurrency: 100,
		}.Create(),
		Summoners: structures.QueueSettings{
			Concurrency: 100,
		}.Create(),
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
		s.Games.Offer(strconv.Itoa(g.GameId))
		for _, p := range g.Participants {
			s.Games.Offer(p.SummonerName)
		}
	}
	return nil
}
