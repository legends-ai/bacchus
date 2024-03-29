package rank

import (
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/asunaio/bacchus/config"
	"github.com/asunaio/bacchus/db"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/golang/protobuf/ptypes"
)

// LookupService looks things up.
type LookupService struct {
	Batcher  *Batcher          `inject:"t"`
	Config   *config.AppConfig `inject:"t"`
	Logger   *logrus.Logger    `inject:"t"`
	Rankings *db.RankingsDAO   `inject:"t"`
}

// Lookup looks up the given ids for a time and returns a rank.
func (ls *LookupService) Lookup(ids []*apb.SummonerId) (map[*apb.SummonerId]*apb.Ranking, error) {
	var err error
	var mu sync.Mutex
	var wg sync.WaitGroup

	ret := map[*apb.SummonerId]*apb.Ranking{}
	wg.Add(len(ids))
	for _, id := range ids {
		// Asynchronously look up all summoners
		go func(id *apb.SummonerId) {
			defer wg.Done()
			ranking, err := ls.lookup(id)
			if err != nil {
				ls.Logger.Errorf("Error looking up ranking: %v", err)
				return
			}
			if ranking == nil {
				return
			}
			mu.Lock()
			ret[id] = ranking
			mu.Unlock()
		}(id)
	}
	wg.Wait()

	// check for one failure
	if err != nil {
		return nil, fmt.Errorf("could not lookup ranking: %v", err)
	}

	return ret, nil
}

func (ls *LookupService) lookup(id *apb.SummonerId) (*apb.Ranking, error) {
	// check cassandra cache
	ranking, err := ls.lookupCassandra(id)
	if err != nil {
		return nil, err
	}
	if ranking != nil {
		return ranking, nil
	}

	// not in cassandra, do api lookup
	ranking, err = ls.Batcher.Lookup(id)
	if err != nil {
		return nil, fmt.Errorf("could not lookup ranking for %v: %v", id, err)
	}

	if ranking == nil {
		// missing ranking from riot
		return nil, nil
	}

	if err = ls.Rankings.Insert(ranking); err != nil {
		return nil, fmt.Errorf("error inserting ranking: %v", err)
	}

	return ranking, nil
}

// lookupCassandra looks up the summoner rank in Cassandra.
// Returns the rank if it exists, whether the rank is already in Cassandra,
// and an error if it exists.
func (ls *LookupService) lookupCassandra(id *apb.SummonerId) (*apb.Ranking, error) {
	// check cassandra cache
	res, err := ls.Rankings.Get(id)
	if err != nil {
		return nil, fmt.Errorf("could not lookup Cassandra: %v", err)
	}
	if res == nil {
		return nil, nil
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
	return res, nil
}
