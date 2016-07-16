package models

import (
	"fmt"
	"strconv"
	"strings"
)

// SummonerID identifies a summoner.
type SummonerID struct {
	Region string
	ID     int
}

// SummonerIDFromString returns a new SummonerID from a string.
func SummonerIDFromString(id string) (SummonerID, error) {
	parts := strings.Split(id, "/")
	i, err := strconv.Atoi(parts[1])
	if err != nil {
		return SummonerID{}, err
	}
	return SummonerID{parts[0], i}, nil
}

// String returns a string representation of this ID.
func (id SummonerID) String() string {
	return fmt.Sprintf("%s/%d", id.Region, id.ID)
}
