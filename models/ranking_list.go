package models

import (
	"sort"
	"time"
)

// RankingList is a list of rankings.
type RankingList struct {
	Rankings []*Ranking
	s        bool // sorted?
}

// NewRankingList creates a new ranking list.
func NewRankingList(rankings []*Ranking) *RankingList {
	return &RankingList{Rankings: rankings}
}

// Len returns the length of the list.
func (r *RankingList) Len() int {
	return len(r.Rankings)
}

// Less reports which is less.
func (r *RankingList) Less(i, j int) bool {
	return r.Rankings[i].Time.Before(r.Rankings[j].Time)
}

// Swap swaps two elements.
func (r *RankingList) Swap(i, j int) {
	r.Rankings[i], r.Rankings[j] = r.Rankings[j], r.Rankings[i]
}

// AtTime returns the ranking active at the given time.
func (r *RankingList) AtTime(t time.Time) *Ranking {
	r.mustBeSorted()
	if r.Len() == 0 {
		return nil
	}
	if t.Before(r.Rankings[0].Time) {
		return r.Rankings[0]
	}
	for i := 0; i < r.Len()-1; i++ {
		if (r.Rankings[i].Time.Before(t) || r.Rankings[i].Time.Equal(t)) && r.Rankings[i+1].Time.After(t) {
			return r.Rankings[i]
		}
	}
	return r.Latest()
}

// Latest gets the latest ranking.
func (r *RankingList) Latest() *Ranking {
	r.mustBeSorted()
	if r.Len() == 0 {
		return nil
	}
	return r.Rankings[len(r.Rankings)-1]
}

func (r *RankingList) mustBeSorted() {
	if !r.s {
		sort.Sort(r)
		r.s = true
	}
}
