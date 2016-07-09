package models

import "time"

// Ranking is a rank and a time fetched.
type Ranking struct {
	Time time.Time
	Rank Rank
}

// RankingList is a list of rankings.
type RankingList struct {
	Rankings []*Ranking
}

// Latest gets the latest ranking.
func (r *RankingList) Latest() *Ranking {
	// TODO(igm): implement
	return nil
}
