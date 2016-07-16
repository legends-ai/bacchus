package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRankingAtTime(t *testing.T) {
	id := SummonerID{}
	for i, test := range []struct {
		Rankings []*Ranking
		Time     time.Time
		Expected int
	}{
		// Middle should use middle
		{
			Rankings: []*Ranking{
				{id, parseTime("2016-01-01"), Rank{10, 10}},
				{id, parseTime("2016-02-01"), Rank{10, 20}},
				{id, parseTime("2016-03-01"), Rank{10, 40}},
			},
			Time:     parseTime("2016-02-10"),
			Expected: 1,
		},
		// After should use the last ranking
		{
			Rankings: []*Ranking{
				{id, parseTime("2016-01-01"), Rank{10, 10}},
				{id, parseTime("2016-02-01"), Rank{10, 20}},
				{id, parseTime("2016-03-01"), Rank{10, 40}},
			},
			Time:     parseTime("2017-02-10"),
			Expected: 2,
		},
		// Before should use the first ranking
		{
			Rankings: []*Ranking{
				{id, parseTime("2016-01-01"), Rank{10, 10}},
				{id, parseTime("2016-02-01"), Rank{10, 20}},
				{id, parseTime("2016-03-01"), Rank{10, 40}},
			},
			Time:     parseTime("2011-02-10"),
			Expected: 0,
		},
		// Before first ranking
		{
			Rankings: []*Ranking{
				{id, parseTime("2016-01-01"), Rank{10, 10}},
				{id, parseTime("2016-02-01"), Rank{10, 20}},
				{id, parseTime("2016-03-01"), Rank{10, 40}},
			},
			Time:     parseTime("2016-01-10"),
			Expected: 0,
		},
		// Sort properly
		{
			Rankings: []*Ranking{
				{id, parseTime("2016-03-01"), Rank{10, 40}},
				{id, parseTime("2016-01-01"), Rank{10, 10}},
				{id, parseTime("2016-02-01"), Rank{10, 20}},
			},
			Time:     parseTime("2016-02-10"),
			Expected: 2,
		},
		// Empty list
		{
			Rankings: []*Ranking{},
			Time:     parseTime("2016-02-10"),
			Expected: -1,
		},
	} {
		// copy since we mutate the slice
		var rankings []*Ranking
		for _, el := range test.Rankings {
			rankings = append(rankings, el)
		}
		ranking := NewRankingList(rankings).AtTime(test.Time)
		if test.Expected != -1 {
			assert.Equal(t, test.Rankings[test.Expected], ranking, fmt.Sprintf("Failed test %v -- Expected %v, got %v", i, test.Expected, ranking))
		} else {
			assert.Nil(t, ranking)
		}
	}
}

func parseTime(t string) time.Time {
	ret, _ := time.Parse("2006-01-02", t)
	return ret
}
