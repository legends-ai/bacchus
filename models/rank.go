package models

import (
	"fmt"

	"github.com/simplyianm/bacchus/riotclient"
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

// RankFromNumber returns a Rank from a number.
func RankFromNumber(n uint32) Rank {
	return Rank{
		Division: uint16(n >> 16),
		Tier:     uint16(n & 0xffff),
	}
}

// ParseRank parses a tier and division to return a Rank.
func ParseRank(tier, division string) (*Rank, error) {
	di := parseDivision(division)
	ti := parseTier(tier)
	if di == 0 {
		return nil, fmt.Errorf("Invalid division %s", division)
	}
	if ti == 0 {
		return nil, fmt.Errorf("Invalid tier %s", tier)
	}
	return &Rank{di, ti}, nil
}

func parseTier(s string) uint16 {
	switch s {
	case riotclient.TierChallenger:
		return 70
	case riotclient.TierMaster:
		return 60
	case riotclient.TierDiamond:
		return 50
	case riotclient.TierPlatinum:
		return 40
	case riotclient.TierGold:
		return 30
	case riotclient.TierSilver:
		return 20
	case riotclient.TierBronze:
		return 10
	default:
		return 0
	}
}

func parseDivision(s string) uint16 {
	switch s {
	case "I":
		return 50
	case "II":
		return 40
	case "III":
		return 30
	case "IV":
		return 20
	case "V":
		return 10
	default:
		return 0
	}
}
