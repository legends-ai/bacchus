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
	List   []string
	Set    string
	c      chan *apb.MatchId
	mx     sync.RWMutex
}

func NewMatchQueue() *MatchQueue {
	return &MatchQueue{
		List: []string{},
		Set:  "SMATCH",
		c:    make(chan *apb.MatchId, 10),
	}
}

func (q *MatchQueue) Start() {
	for {
		if len(q.List) < 1 {
			continue
		}
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

		if set, err := q.Redis.SMembers(q.Set).Result(); err != nil {
			q.Logger.Warnf("SMEMBERS %v failed: %v", q.Set, err)
			return
		} else {
			sort.Sort(MatchQueueList(set))
			q.mx.Lock()
			q.List = set
			q.mx.Unlock()
		}
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
