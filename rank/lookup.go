package rank

import (
	"fmt"
	"strconv"
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
	Logger *logrus.Logger         `inject:"t"`
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
		go func(id models.SummonerID) {
			rank, err := ls.lookup(id, t)
			if err != nil {
				ls.Logger.Errorf("Error looking up rank: %v", err)
			}
			mu.Lock()
			ret[id] = *rank
			mu.Unlock()
		}(id)
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
	rank, err := ls.lookupCassandra(id, t)
	if err != nil {
		return nil, err
	}
	if rank != nil {
		return rank, nil
	}
	// not in cassandra, do api lookup
	ls.Logger.Infof("Expired rank for %s, performing API lookup", id.String())
	r := ls.Riot.Region(id.Region)
	// TODO(igm): batch id lookups. we can fit a lot of these in a URI.
	res, err := r.League([]string{strconv.Itoa(id.ID)})
	if err != nil {
		return nil, err
	}

	dtos := res[strconv.Itoa(id.ID)]
	var dto *riotclient.LeagueDto
	for _, x := range dtos {
		if x.Queue == riotclient.QueueSolo5x5 {
			dto = x
			break
		}
	}
	if dto == nil {
		// unranked
		return nil, nil
	}

	tier := dto.Tier
	var entry *riotclient.LeagueEntryDto
	for _, x := range dto.Entries {
		if x.PlayerOrTeamID == strconv.Itoa(id.ID) {
			entry = x
			break
		}
	}
	if entry == nil {
		// should not happen
		return nil, fmt.Errorf("No summoner %d for league %s of %s", id.ID, dto.Name, dto.Tier)
	}

	rank, err = models.ParseRank(entry.Division, tier)
	if err != nil {
		return nil, fmt.Errorf("Invalid rank: %v", err)
	}

	// asynchronously update cassandra
	go ls.updateCassandra(id, t, *rank)
	return rank, nil
}

// lookupCassandra looks up the summoner rank in Cassandra.
func (ls *LookupService) lookupCassandra(id models.SummonerID, t time.Time) (*models.Rank, error) {
	// check cassandra cache
	res, err := ls.Athena.Rankings(id)
	if err != nil {
		return nil, fmt.Errorf("could not lookup Cassandra: %v", err)
	}
	ranking := res.AtTime(t)
	if ranking == nil {
		return nil, nil
	}
	if ranking != res.Latest() || time.Now().Sub(ranking.Time) < ls.Config.RankExpiry {
		return &ranking.Rank, nil
	}
	return nil, nil
}

func (ls *LookupService) updateCassandra(id models.SummonerID, t time.Time, rank models.Rank) {
	// TODO implement
}
