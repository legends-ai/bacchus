package rank

import (
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/models"
	"github.com/simplyianm/bacchus/riotclient"
)

// Number of players to batch lookups at once
const batchSize = 20

// subscription is a subscription to later complete.
type subscription struct {
	id models.SummonerID
	c  chan riotclient.LeagueResponse
	e  chan error
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
			for _, s := range subs {
				s.e <- err
				close(s.c)
				close(s.e)
			}
			return
		}

		// return results
		for _, s := range subs {
			s.c <- res
			close(s.c)
			close(s.e)
		}
	}
}

// subscribe subscribes to the response generated from looking up an id.
func (b *batchRegion) subscribe(id models.SummonerID) (riotclient.LeagueResponse, error) {
	sub := &subscription{
		id: id,
		c:  make(chan riotclient.LeagueResponse),
		e:  make(chan error),
	}
	b.subs <- sub
	select {
	case res := <-sub.c:
		return res, nil
	case err := <-sub.e:
		return nil, err
	}
}

// Batcher batches ranking lookups from Riot.
type Batcher struct {
	Riot   *riotclient.RiotClient `inject:"t"`
	Logger *logrus.Logger         `inject:"t"`

	batchers map[string]*batchRegion
	mu       sync.Mutex
}

// NewBatcher constructs a new batcher for rank lookups.
func NewBatcher() *Batcher {
	return &Batcher{
		batchers: map[string]*batchRegion{},
	}
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
func (b *Batcher) Lookup(id models.SummonerID) ([]*riotclient.LeagueDto, error) {
	res, err := b.Region(id.Region).subscribe(id)
	if err != nil {
		return nil, err
	}
	return res[strconv.Itoa(id.ID)], nil
}
