package queue

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v4"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
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
	}
}

func (q *SummonerQueue) Start() {
	for {
		r, err := q.Redis.BLPop(0, q.List...).Result()
		if err != nil {
			q.Logger.Warnf("BLPOP %v failed: %v", q.List, err)
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
	list := fmt.Sprintf("%#x", ctx.Rank.Tier)
	summoner := in.String()

	_, err := q.Redis.RPush(list, summoner).Result()
	if err != nil {
		q.Logger.Warnf("RPUSH %v to %v failed: %v", summoner, list, err)
	}
}

func (q *SummonerQueue) Poll() *apb.SummonerId {
	return <-q.c
}
