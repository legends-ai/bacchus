package queue

import (
	"bytes"
	"encoding/gob"

	"gopkg.in/redis.v4"

	"github.com/asunaio/bacchus/config"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
)

type MatchQueue struct {
	c         *redis.Client `inject:"t"`
	MatchList string
}

func init() {
	gob.Register(apb.MatchId{})
}

func NewMatchQueue() *MatchQueue {
	return &MatchQueue{
		c: redis.NewClient(&redis.Options{
			Addr:     config.Fetch().RedisHost,
			Password: "",
			DB:       0,
		}),
		MatchList: "MATCH",
	}
}

func (q *MatchQueue) Add(in *apb.MatchId, ctx *apb.CharonMatchListResponse_MatchInfo) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(*in); err != nil {
		return
	}
	q.c.RPush(q.MatchList, b.Bytes())
}

func (q *MatchQueue) Poll() *apb.MatchId {
	r, err := q.c.BLPop(0, q.MatchList).Result()
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
