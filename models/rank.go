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

// RankFromNumber returns a Rank from a number.
func RankFromNumber(n uint32) Rank {
	return Rank{
		Division: uint16(n >> 16),
		Tier:     uint16(n & 0xffff),
	}
}
