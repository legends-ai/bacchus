package queue

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v4"

	"github.com/asunaio/bacchus/config"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
)

type MatchQueue struct {
	c    *redis.Client
	List []string
}

func NewMatchQueue() *MatchQueue {
	return &MatchQueue{
		c: redis.NewClient(&redis.Options{
			Addr:     config.Fetch().RedisHost,
			Password: "",
			DB:       0,
		}),
		List: []string{"MATCH"},
	}
}

func (q *MatchQueue) Add(in interface{}, ctx interface{}) {
	q.c.RPush(q.List[0], in.(*apb.MatchId).String())
}

func (q *MatchQueue) Poll() interface{} {
	r, err := q.c.BLPop(0, q.List...).Result()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	data := &apb.MatchId{}
	if err := proto.UnmarshalText(r[1], data); err != nil {
		return nil
	}
	return data
}
