package models

import "time"

// RankingUDT is the UDT for a ranking.
type RankingUDT struct {
	Time time.Time `cql:"time"`
	Rank uint32    `cql:"rank"`
}

// ToRanking converts this UDT to a ranking.
func (r RankingUDT) ToRanking() *Ranking {
	return &Ranking{
		Time: r.Time,
		Rank: RankFromNumber(r.Rank),
	}
}
