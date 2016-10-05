package queue

import (
	"bytes"
	"encoding/gob"
	"fmt"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"gopkg.in/redis.v4"
)

type SummonerQueue struct {
	c             *redis.Client `inject:"t"`
	SummonerLists []string
}

func init() {
	gob.Register(apb.SummonerId{})
}

func NewSummonerQueue() *SummonerQueue {
	/*return &SummonerQueue{
	c: redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}),*/
	return &SummonerQueue{
		SummonerLists: []string{
			"0x70", "0x60", "0x50", "0x40", "0x30", "0x20", "0x10",
		},
	}
}

func (q *SummonerQueue) Add(in *apb.SummonerId, ctx *apb.Ranking) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	if err := e.Encode(*in); err != nil {
		return
	}
	fmt.Println(q.c)
	q.c.RPush(fmt.Sprintf("%#x", ctx.Rank.Tier), b.Bytes())
}

func (q *SummonerQueue) Poll() *apb.SummonerId {
	r, err := q.c.BLPop(0, q.SummonerLists...).Result()
	if err != nil {
		return nil
	}
	data := &apb.SummonerId{}
	b := bytes.Buffer{}
	b.Write([]byte(r[1]))
	d := gob.NewDecoder(&b)
	err = d.Decode(data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}
