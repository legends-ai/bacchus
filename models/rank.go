package models

import (
	"fmt"
	"sort"

	apb "github.com/asunaio/bacchus/gen-go/asuna"
)

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

type rankAscending []*apb.Rank

func (r rankAscending) Swap(a, b int) {
	r[a], r[b] = r[b], r[a]
}

func (r rankAscending) Less(a, b int) bool {
	return RankToNumber(r[a]) < RankToNumber(r[b])
}

func (r rankAscending) Len() int {
	return len(r)
}

// RankToNumber returns a numerical representation of rank that can be sorted.
func RankToNumber(r *apb.Rank) uint32 {
	return r.Tier<<16 | r.Division
}

// RankFromNumber returns a Rank from a number.
func RankFromNumber(n uint32) *apb.Rank {
	return &apb.Rank{
		Tier:     n >> 16,
		Division: n & 0xffff,
	}
}

// MedianRank calculates the median rank in the list.
func MedianRank(res []*apb.Rank) *apb.Rank {
	l := len(res)
	if l == 0 {
		return nil
	}

	sort.Sort(rankAscending(res))
	if l%2 == 0 {
		return res[l/2-1]
	} else {
		return res[l/2]
	}
}

// MinRank gets the minimum rank out of the given ranks
func MinRank(res []*apb.Rank) *apb.Rank {
	min := &apb.Rank{1<<16 - 1, 1<<16 - 1}
	for _, rank := range res {
		if RankOver(rank, min) {
			continue
		}
		min = rank
	}
	return min
}

// MaxRank gets the maximum rank out of the given ranks
func MaxRank(res []*apb.Rank) *apb.Rank {
	max := &apb.Rank{1<<16 - 1, 1<<16 - 1}
	for _, rank := range res {
		if RankOver(rank, max) {
			max = rank
		}
	}
	return max
}

func RankOver(a *apb.Rank, b *apb.Rank) bool {
	return RankToNumber(a) > RankToNumber(b)
}

// ParseRank parses a tier and division to return a Rank.
func ParseRank(tier, division string) (*apb.Rank, error) {
	ti := parseTier(tier)
	if ti == 0 {
		return nil, fmt.Errorf("invalid tier %s", tier)
	}
	di := parseDivision(division)
	if di == 0 {
		return nil, fmt.Errorf("invalid division %s", division)
	}
	return &apb.Rank{
		Tier:     ti,
		Division: di,
	}, nil
}

func parseTier(s string) uint32 {
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

func parseDivision(s string) uint32 {
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
