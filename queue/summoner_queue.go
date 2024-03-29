package queue

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v5"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/models"
)

type SummonerQueue struct {
	Logger *logrus.Logger `inject:"t"`
	Redis  *redis.Client  `inject:"t"`
	List   []string
	c      chan *apb.SummonerId
}

func NewSummonerQueue() *SummonerQueue {
	return &SummonerQueue{
		List: []string{
			"0x70", "0x60", "0x50", "0x40", "0x30", "0x20", "0x10",
		},
		c: make(chan *apb.SummonerId, 10),
	}
}

func (q *SummonerQueue) Start() {
	for {
		r, err := q.Redis.BLPop(0, q.List...).Result()
		if err != nil {
			q.Logger.Warnf("BLPOP %v failed: %v", q.List, err)
			continue
		}

		set := fmt.Sprintf("S%s", r[0])
		_, err = q.Redis.SRem(set, r[1]).Result()
		for i := 1; err != nil && i <= 3; i++ {
			q.Logger.Warnf("SREM %v from %v failed: %v", set, r[0], err)
			q.Logger.Warnf("Retry %v | Retries left: %v", i, 3-i)
			_, err = q.Redis.SRem(set, r[1]).Result()
		}

		if err != nil {
			q.Logger.Warnf("SREM %v from %v failed 3 times: %v", set, r[0], err)
			q.Logger.Warnf("Moving on - %v will never be parsed again...", r[0])
			continue
		}

		var data apb.SummonerId
		if err := proto.UnmarshalText(r[1], &data); err != nil {
			q.Logger.Warnf("UnmarshalText %v failed: %v", r[1], err)
			continue
		}
		q.c <- &data
	}
}

func (q *SummonerQueue) Add(in *apb.SummonerId, ctx *apb.Ranking) {
	// Ignore the ranking if there is none.
	if len(ctx.Ranks) == 0 {
		return
	}

	// Extract the apb.Rank from each QueueRank
	var theRanks []*apb.Rank
	for _, rank := range ctx.Ranks {
		theRanks = append(theRanks, rank.Rank)
	}

	list := fmt.Sprintf("%#x", models.MaxRank(theRanks).Tier)
	set := fmt.Sprintf("S%s", list)
	summoner := in.String()

	if exists, err := q.Redis.SIsMember(set, summoner).Result(); err != nil {
		q.Logger.Warnf("SISMEMBER %v in %v failed: %v", summoner, set, err)
	} else if exists {
		return
	}

	if llen, err := q.Redis.LLen(list).Result(); err != nil {
		q.Logger.Warnf("LLEN %v failed (skipping this summoner): %v", list, err)
		return
	} else if llen >= 1000000 {
		return
	}

	if _, err := q.Redis.SAdd(set, summoner).Result(); err != nil {
		q.Logger.Warnf("SADD %v to %v failed: %v", summoner, set, err)
		return
	}

	if _, err := q.Redis.RPush(list, summoner).Result(); err != nil {
		q.Logger.Warnf("RPUSH %v to %v failed: %v", summoner, list, err)
	}
}

func (q *SummonerQueue) Poll() *apb.SummonerId {
	return <-q.c
}
