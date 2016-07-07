package rank

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/riotclient"
)

// Rank represents a rank.
type Rank struct {
	Division int
	Tier     int
}

// LookupService looks things up.
type LookupService struct {
	Riot   *riotclient.RiotClient `inject:"t"`
	Logger logrus.Logger          `inject:"t"`
}

// Lookup looks up the given ids and returns a rank.
func (ls *LookupService) Lookup(ids []db.SummonerID) map[db.SummonerID]Rank {
	var mu sync.Mutex
	var wg sync.WaitGroup
	ret := map[db.SummonerID]Rank{}
	wg.Add(len(ids))
	for _, id := range ids {
		// Asynchronously look up all summoners
		go func() {
			rank := lookup(id)
			mu.Lock()
			ret[id] = rank
			mu.Unlock()
		}()
	}
	return ret
}

func lookup(id db.SummonerID) Rank {
	// TODO(simplyianm): implement
	return Rank{}
}
