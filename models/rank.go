package models

import "fmt"

const (
	TierChallenger = "CHALLENGER"
	TierMaster     = "MASTER"
	TierDiamond    = "DIAMOND"
	TierPlatinum   = "PLATINUM"
	TierGold       = "GOLD"
	TierSilver     = "SILVER"
	TierBronze     = "BRONZE"

	DivisionI   = "I"
	DivisionII  = "II"
	DivisionIII = "III"
	DivisionIV  = "IV"
	DivisionV   = "V"
)

// Rank represents a rank.
type Rank struct {
	Tier     uint16
	Division uint16
}

// ToNumber returns a numerical representation of rank that can be sorted.
func (r Rank) ToNumber() uint32 {
	return uint32(r.Tier)<<16 | uint32(r.Division)
}

// Over checks if the rank is >= the given rank
func (r Rank) Over(o Rank) bool {
	return r.ToNumber() >= o.ToNumber()
}

// RankFromNumber returns a Rank from a number.
func RankFromNumber(n uint32) Rank {
	return Rank{
		Division: uint16(n & 0xffff),
		Tier:     uint16(n >> 16),
	}
}

// MinRank gets the minimum rank out of the given ranks
func MinRank(res []Rank) Rank {
	min := Rank{1<<16 - 1, 1<<16 - 1}
	for _, rank := range res {
		if !rank.Over(min) {
			min = rank
		}
	}
	return min
}

// ParseRank parses a tier and division to return a Rank.
func ParseRank(tier, division string) (*Rank, error) {
	ti := parseTier(tier)
	if ti == 0 {
		return nil, fmt.Errorf("invalid tier %s", tier)
	}
	di := parseDivision(division)
	if di == 0 {
		return nil, fmt.Errorf("invalid division %s", division)
	}
	return &Rank{
		Tier:     ti,
		Division: di,
	}, nil
}

func parseTier(s string) uint16 {
	switch s {
	case TierChallenger:
		return 0x70
	case TierMaster:
		return 0x60
	case TierDiamond:
		return 0x50
	case TierPlatinum:
		return 0x40
	case TierGold:
		return 0x30
	case TierSilver:
		return 0x20
	case TierBronze:
		return 0x10
	default:
		return 0
	}
}

func parseDivision(s string) uint16 {
	switch s {
	case DivisionI:
		return 0x50
	case DivisionII:
		return 0x40
	case DivisionIII:
		return 0x30
	case DivisionIV:
		return 0x20
	case DivisionV:
		return 0x10
	default:
		return 0
	}
}
