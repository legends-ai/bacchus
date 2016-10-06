package queue

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v4"

	"github.com/asunaio/bacchus/config"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
)

type SummonerQueue struct {
	c    *redis.Client `inject:"t"`
	List []string
}

func NewSummonerQueue() *SummonerQueue {
	return &SummonerQueue{
		c: redis.NewClient(&redis.Options{
			Addr:     config.Fetch().RedisHost,
			Password: "",
			DB:       0,
		}),
		List: []string{
			"0x70", "0x60", "0x50", "0x40", "0x30", "0x20", "0x10",
		},
	}
}

func (q *SummonerQueue) Add(in interface{}, ctx interface{}) {
	list := fmt.Sprintf("%#x", ctx.(*apb.Ranking).Rank.Tier)
	q.c.RPush(list, in.(*apb.SummonerId).String())
}

func (q *SummonerQueue) Poll() interface{} {
	r, err := q.c.BLPop(0, q.List...).Result()
	if err != nil {
		return nil
	}
	data := &apb.SummonerId{}
	if err := proto.UnmarshalText(r[1], data); err != nil {
		return nil
	}
	return data
}
