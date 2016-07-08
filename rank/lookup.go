package rank

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/riotclient"
)

// LookupService looks things up.
type LookupService struct {
	Riot   *riotclient.RiotClient `inject:"t"`
	Logger logrus.Logger          `inject:"t"`
}

// Lookup looks up the given ids and returns a rank.
func (ls *LookupService) Lookup(ids []db.SummonerID) map[db.SummonerID]db.Rank {
	var mu sync.Mutex
	var wg sync.WaitGroup
	ret := map[db.SummonerID]db.Rank{}
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

// MinRank gets the minimum rank of the given summoners.
func (ls *LookupService) MinRank(ids []db.SummonerID) db.Rank {
	res := ls.Lookup(ids)
	min := db.Rank{1<<16 - 1, 1<<16 - 1}
	for _, rank := range res {
		if rank.Tier > min.Tier {
			continue
		}
		if rank.Division > min.Division {
			continue
		}
		min = rank
	}
	return min
}

func lookup(id db.SummonerID) db.Rank {
	// TODO(simplyianm): implement
	return db.Rank{}
}
