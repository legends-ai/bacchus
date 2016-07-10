package rank

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/config"
	"github.com/simplyianm/bacchus/db"
	"github.com/simplyianm/bacchus/models"
	"github.com/simplyianm/bacchus/riotclient"
)

// LookupService looks things up.
type LookupService struct {
	Riot   *riotclient.RiotClient `inject:"t"`
	Logger logrus.Logger          `inject:"t"`
	Athena *db.Athena             `inject:"t"`
	Config *config.AppConfig      `inject:"t"`
}

// Lookup looks up the given ids for a time and returns a rank.
func (ls *LookupService) Lookup(ids []models.SummonerID, t time.Time) map[models.SummonerID]models.Rank {
	var mu sync.Mutex
	var wg sync.WaitGroup
	ret := map[models.SummonerID]models.Rank{}
	wg.Add(len(ids))
	for _, id := range ids {
		// Asynchronously look up all summoners
		go func() {
			rank, err := ls.lookup(id, t)
			if err != nil {
				ls.Logger.Errorf("Error looking up rank: %v", err)
			}
			mu.Lock()
			ret[id] = *rank
			mu.Unlock()
		}()
	}
	return ret
}

// MinRank gets the minimum rank of the given summoners.
func (ls *LookupService) MinRank(ids []models.SummonerID, t time.Time) models.Rank {
	res := ls.Lookup(ids, t)
	min := models.Rank{1<<16 - 1, 1<<16 - 1}
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

func (ls *LookupService) lookup(id models.SummonerID, t time.Time) (*models.Rank, error) {
	// check cassandra cache
	res, err := ls.Athena.Rankings(id)
	if err != nil {
		return nil, err
	}
	ranking := res.AtTime(t)
	if ranking == res.Latest() {

	}
	if ranking != nil {
		// TODO(igm): check timestamp on result
		return &ranking.Rank, nil
	}
	// lookup true rank and update
	return &models.Rank{}, nil
}
