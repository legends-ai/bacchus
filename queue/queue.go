package queue

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"gopkg.in/redis.v4"

	"github.com/asunaio/bacchus/config"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
)

// Queue is a priority queue.
type Queue interface {
	// Add adds an element to the queue.
	// Context is an arbitrary value which may influence the priority of the element within the queue.
	// Summoner: in *apb.SummonerId, context *apb.Ranking
	// Match: in *apb.MatchId, context *apb.CharonMatchListResponse_MatchInfo
	Add(in interface{}, context interface{})

	// Poll gets the next element of the queue to process.
	Poll() interface{}
}

type RedisQueue struct {
	c      *redis.Client
	decode interface{}
	List   []string
}

func NewRedisQueue(list []string, decode interface{}) *RedisQueue {
	return &RedisQueue{
		c: redis.NewClient(&redis.Options{
			Addr:     config.Fetch().RedisHost,
			Password: "",
			DB:       0,
		}),
		decode: decode,
		List:   list,
	}
}

func (q *RedisQueue) Add(in interface{}, ctx interface{}) {
	switch in.(type) {
	case *apb.SummonerId:
		data, err := q.serialize(in)
		if err != nil {
			return
		}
		q.c.RPush(fmt.Sprintf("%#x", ctx.(*apb.Ranking).Rank.Tier), data)
	case *apb.MatchId:
		data, err := q.serialize(in)
		if err != nil {
			return
		}
		fmt.Println(data)
		q.c.RPush(q.List[0], data)
	}
}

func (q *RedisQueue) Poll() interface{} {
	r, err := q.c.BLPop(0, q.List...).Result()
	if err != nil {
		return nil
	}
	data, err := q.deserialize(r[1])
	if err != nil {
		return nil
	}
	return data
}

func (q *RedisQueue) serialize(in interface{}) ([]byte, error) {
	var err error
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)

	switch in.(type) {
	case *apb.SummonerId:
		err = e.Encode(*in.(*apb.SummonerId))
	case *apb.MatchId:
		err = e.Encode(*in.(*apb.MatchId))
	}
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (q *RedisQueue) deserialize(r string) (interface{}, error) {
	var err error
	var data interface{}
	b := bytes.Buffer{}
	b.Write([]byte(r))
	d := gob.NewDecoder(&b)

	switch q.decode.(type) {
	case *apb.SummonerId:
		data = &apb.SummonerId{}
	case *apb.MatchId:
		data = &apb.MatchId{}
	}

	err = d.Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
