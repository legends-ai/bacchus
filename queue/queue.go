package queue

import (
	"bytes"
	"encoding/gob"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"gopkg.in/redis.v4"
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

func init() {
	gob.Register(apb.SummonerId{})
	gob.Register(apb.MatchId{})
}

type RedisQueue struct {
	c    *redis.Client
	List []string
}

func NewRedisQueue() *RedisQueue {
	return &RedisQueue{
		c: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
		List: []string{"MATCH"},
	}
}

func (q *RedisQueue) Add(in *apb.MatchId, ctx *apb.CharonMatchListResponse_MatchInfo) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(*in); err != nil {
		return
	}
	q.c.RPush(q.List[1], b.Bytes())
}

func (q *RedisQueue) Poll() *apb.MatchId {
	r, err := q.c.BLPop(0, q.List...).Result()
	if err != nil {
		return nil
	}
	data := &apb.MatchId{}
	b := bytes.Buffer{}
	b.Write([]byte(r[1]))
	d := gob.NewDecoder(&b)
	err = d.Decode(data)
	if err != nil {
		return nil
	}
	return data
}
