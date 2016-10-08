package queue

import (
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v4"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
)

type MatchQueue struct {
	Logger *logrus.Logger `inject:"t"`
	Redis  *redis.Client  `inject:"t"`
	List   []string
}

func NewMatchQueue() *MatchQueue {
	return &MatchQueue{
		List: []string{"MATCH"},
	}
}

func (q *MatchQueue) Add(in *apb.MatchId, ctx *apb.CharonMatchListResponse_MatchInfo) {
	list := q.List[0]
	match := in.String()

	_, err := q.Redis.RPush(list, match).Result()
	if err != nil {
		q.Logger.Warnf("RPUSH %v to %v failed: %v", match, list, err)
	}
}

func (q *MatchQueue) Poll() *apb.MatchId {
	r, err := q.Redis.BLPop(0, q.List...).Result()
	if err != nil {
		q.Logger.Warnf("BLPOP %v failed: %v", q.List, err)
		return nil
	}
	data := &apb.MatchId{}
	if err := proto.UnmarshalText(r[1], data); err != nil {
		q.Logger.Warnf("UnmarshalText %v failed: %v", r[1], err)
		return nil
	}
	return data
}
