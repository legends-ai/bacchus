package models

import "time"

// Ranking is a rank and a time fetched.
type Ranking struct {
	ID   SummonerID
	Time time.Time
	Rank Rank
}
