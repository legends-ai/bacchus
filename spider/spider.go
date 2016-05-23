package spider

import (
	"fmt"

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
	r, err := s.Riot.FeaturedGames()
	if err != nil {
		return nil, fmt.Errorf("Could not get featured games: %v", err)
	}
	return s, nil
}
