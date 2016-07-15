package models

import "time"

// Ranking is a rank and a time fetched.
type Ranking struct {
	Time time.Time
	Rank Rank
}

// UDTSet returns a UDT set type for this Ranking.
func (r Ranking) UDTSet() []RankingUDT {
	return []RankingUDT{{r.Time, r.Rank.ToNumber()}}
}
