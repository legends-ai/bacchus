package queue

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v4"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
)

type MatchQueueList []string

func (s MatchQueueList) Len() int {
	return len(s)
}

func (s MatchQueueList) Less(i, j int) bool {
	var imonth, iyear, jmonth, jyear int
	fmt.Sscanf(s[i], "MATCH:%d %d", &imonth, &iyear)
	fmt.Sscanf(s[j], "MATCH:%d %d", &jmonth, &jyear)

	itime := time.Date(iyear, time.Month(imonth), 1, 0, 0, 0, 0, time.UTC)
	jtime := time.Date(jyear, time.Month(jmonth), 1, 0, 0, 0, 0, time.UTC)

	return itime.After(jtime)
}

func (s MatchQueueList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type MatchQueue struct {
	Logger *logrus.Logger `inject:"t"`
	Redis  *redis.Client  `inject:"t"`

	// List is the list of match queues. This is used to BLPOP from the most important queue.
	List []string
	Set  string
	c    chan *apb.MatchId
	mx   sync.RWMutex
}

func NewMatchQueue() *MatchQueue {
	return &MatchQueue{
		List: []string{},
		Set:  "SMATCH",
		c:    make(chan *apb.MatchId, 10),
	}
}

func (q *MatchQueue) Start() {
	// refresh queues on startup in case they already exist in redis
	if err := q.refreshQueues(); err != nil {
		q.Logger.Warnf("Could not refresh match queues: %v", err)
		return
	}

	for {
		if len(q.List) < 1 {
			// We sleep so we don't do a 1 Infinite Loop (thanks pradyuman)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// Pop from the most important queue
		r, err := q.Redis.BLPop(0, q.List...).Result()
		if err != nil {
			q.Logger.Warnf("BLPOP %v failed: %v", q.List, err)
			continue
		}
		var data apb.MatchId
		if err := proto.UnmarshalText(r[1], &data); err != nil {
			q.Logger.Warnf("UnmarshalText %v failed: %v", r[1], err)
			continue
		}
		q.c <- &data
	}
}

func (q *MatchQueue) Add(in *apb.MatchId, ctx *apb.CharonRpc_MatchListResponse_MatchInfo) {
	t := time.Unix(ctx.Timestamp.Seconds, int64(ctx.Timestamp.Nanos))
	list := fmt.Sprintf("MATCH:%d %d", t.Month(), t.Year())

	if exists, err := q.Redis.SIsMember(q.Set, list).Result(); err != nil {
		q.Logger.Warnf("SISMEMBER %v in %v failed: %v", list, q.Set, err)
	} else if !exists {
		if _, err := q.Redis.SAdd(q.Set, list).Result(); err != nil {
			q.Logger.Warnf("SADD %v to %v failed: %v", list, q.Set, err)
			return
		}

		if err := q.refreshQueues(); err != nil {
			q.Logger.Warnf("Could not refresh match queues: %v", err)
			return
		}
	}

	if llen, err := q.Redis.LLen(list).Result(); err != nil {
		q.Logger.Warnf("LLEN %v failed (skipping this match): %v", list, err)
		return
	} else if llen >= 1000000 {
		return
	}

	match := in.String()
	q.mx.RLock()
	if _, err := q.Redis.RPush(list, match).Result(); err != nil {
		q.Logger.Warnf("RPUSH %v to %v failed: %v", match, list, err)
	}
	q.mx.RUnlock()
}

func (q *MatchQueue) Poll() *apb.MatchId {
	return <-q.c
}

// refreshQueues refreshes the list of match queues in memory.
func (q *MatchQueue) refreshQueues() error {
	queues, err := q.Redis.Keys("MATCH:*").Result()
	if err != nil {
		return err
	}
	sort.Sort(MatchQueueList(queues))
	q.mx.Lock()
	q.List = queues
	q.mx.Unlock()
	return nil
}
