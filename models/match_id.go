package models

import "fmt"

// MatchID identifies a match.
type MatchID struct {
	Region string
	ID     int
}

// String returns a string representation of this ID.
func (id MatchID) String() string {
	return fmt.Sprintf("%s/%d", id.Region, id.ID)
}
