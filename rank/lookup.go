package rank

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/asunaio/bacchus/config"
	"github.com/asunaio/bacchus/db"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
	"github.com/asunaio/bacchus/riot"
	"github.com/golang/protobuf/ptypes"
)

// LookupService looks things up.
type LookupService struct {
	Riot     *riot.Client      `inject:"t"`
	Logger   *logrus.Logger    `inject:"t"`
	Config   *config.AppConfig `inject:"t"`
	Rankings *db.RankingsDAO   `inject:"t"`
	Batcher  *Batcher          `inject:"t"`
}

// Lookup looks up the given ids for a time and returns a rank.
func (ls *LookupService) Lookup(ids []*apb.SummonerId) (map[*apb.SummonerId]*apb.Rank, error) {
	var err error
	var mu sync.Mutex
	var wg sync.WaitGroup

	ret := map[*apb.SummonerId]*apb.Rank{}
	wg.Add(len(ids))
	for _, id := range ids {
		// Asynchronously look up all summoners
		go func(id *apb.SummonerId) {
			defer wg.Done()
			rank, err := ls.lookup(id)
			if err != nil {
				ls.Logger.Errorf("Error looking up rank: %v", err)
				return
			}
			if rank == nil {
				return
			}
			mu.Lock()
			ret[id] = rank
			mu.Unlock()
		}(id)
	}
	wg.Wait()

	// check for one failure
	if err != nil {
		return nil, fmt.Errorf("could not lookup rank: %v", err)
	}

	return ret, nil
}

func (ls *LookupService) lookup(id *apb.SummonerId) (*apb.Rank, error) {
	// check cassandra cache
	rank, err := ls.lookupCassandra(id)
	if err != nil {
		return nil, err
	}
	if rank != nil {
		return rank, nil
	}
	// not in cassandra, do api lookup
	ls.Logger.Infof("Expired rank for %s, performing API lookup", id.String())
	dtos, err := ls.Batcher.Lookup(id)

	var dto *riot.LeagueDto
	for _, x := range dtos {
		if x.Queue == riot.QueueSolo5x5 {
			dto = x
			break
		}
	}
	if dto == nil {
		// unranked
		return nil, nil
	}

	// Find player ranking
	tier := dto.Tier
	var entry *riot.LeagueEntryDto
	for _, x := range dto.Entries {
		if x.PlayerOrTeamID == strconv.Itoa(int(id.Id)) {
			entry = x
			break
		}
	}
	if entry == nil {
		// should not happen
		return nil, fmt.Errorf("no summoner %d for league %s of %s", id.Id, dto.Name, dto.Tier)
	}
	rank, err = models.ParseRank(tier, entry.Division)
	if err != nil {
		return nil, fmt.Errorf("invalid rank: %v", err)
	}

	now, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		// this err is wtf
		return nil, fmt.Errorf("could not now the time: %v", err)
	}

	ranking := &apb.Ranking{
		Summoner: id,
		Rank:     rank,
		Time:     now,
	}

	ls.Logger.Infof("Found rank of %d: %s %s (%x)", id.Id, tier, entry.Division, models.RankToNumber(rank))

	if err = ls.Rankings.Insert(ranking); err != nil {
		return nil, fmt.Errorf("error inserting ranking: %v", err)
	}

	return rank, nil
}

// lookupCassandra looks up the summoner rank in Cassandra.
// Returns the rank if it exists, whether the rank is already in Cassandra,
// and an error if it exists.
func (ls *LookupService) lookupCassandra(id *apb.SummonerId) (*apb.Rank, error) {
	// check cassandra cache
	res, err := ls.Rankings.Get(id)
	if err != nil {
		return nil, fmt.Errorf("could not lookup Cassandra: %v", err)
	}

	// get time
	t, err := ptypes.Timestamp(res.Time)
	if err != nil {
		return nil, fmt.Errorf("could not parse time: %v", err)
	}

	// Check if rank is expired
	if time.Now().Sub(t) >= ls.Config.RankExpiry {
		return nil, nil
	}

	// ret
	return res.Rank, nil
}
