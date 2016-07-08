package models

import "fmt"

// SummonerID identifies a summoner.
type SummonerID struct {
	Region string
	ID     int
}

// String returns a string representation of this ID.
func (id SummonerID) String() string {
	return fmt.Sprintf("%s/%s", id.Region, id.ID)
}
