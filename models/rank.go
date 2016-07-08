package models

// Rank represents a rank.
type Rank struct {
	Division uint16
	Tier     uint16
}

// ToNumber returns a numerical representation of rank that can be sorted.
func (r Rank) ToNumber() uint32 {
	return uint32(r.Tier)<<16 | uint32(r.Division)
}
