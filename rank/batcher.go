package rank

import (
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/models"
	"github.com/simplyianm/bacchus/riotclient"
	"github.com/simplyianm/riot/config"
)

// Number of players to batch lookups at once
const batchSize = 20

// subscription is a subscription to later complete.
type subscription struct {
	id models.SummonerID
	c  chan riotclient.LeagueResponse
}

// A batch for a region.
type batchRegion struct {
	b *Batcher
	// Region
	r *riotclient.API
	// Channel containing subscriptions
	subs chan *subscription
}

// batch continuously batches ids and performs lookups.
func (b *batchRegion) batch() {
	var subs []*subscription
	for {
		sub := <-b.subs
		subs = append(subs, sub)
		if len(subs) < batchSize {
			continue
		}

		// Perform batched lookup
		var lookup []string
		for _, s := range subs {
			lookup = append(lookup, strconv.Itoa(s.id.ID))
		}
		res, err := b.r.League(lookup)
		if err != nil {
			b.b.Logger.Errorf("Error batching: %v", err)
		}

		// return results
		for _, s := range subs {
			s.c <- res
			close(s.c)
		}
	}
}

// subscribe subscribes to the response generated from looking up an id.
func (b *batchRegion) subscribe(id models.SummonerID) riotclient.LeagueResponse {
	sub := &subscription{
		id: id,
		c:  make(chan riotclient.LeagueResponse),
	}
	b.subs <- sub
	res := <-sub.c
	return res
}

// Batcher batches ranking lookups from Riot.
type Batcher struct {
	Riot   *riotclient.RiotClient `inject:"t"`
	Logger *logrus.Logger         `inject:"t"`
	Config *config.AppConfig      `inject:"t"`

	batchers map[string]*batchRegion
	mu       sync.Mutex
}

// Region is a batching region.
func (b *Batcher) Region(region string) *batchRegion {
	b.mu.Lock()
	defer b.mu.Unlock()
	inst, ok := b.batchers[region]
	if !ok {
		inst = &batchRegion{
			b:    b,
			r:    b.Riot.Region(region),
			subs: make(chan *subscription),
		}
		go inst.batch()
		b.batchers[region] = inst
	}
	return inst
}

// Lookup looks up the id and returns the league response once batched and constructed
func (b *Batcher) Lookup(id models.SummonerID) []*riotclient.LeagueDto {
	res := b.Region(id.Region).subscribe(id)
	return res[strconv.Itoa(id.ID)]
}
