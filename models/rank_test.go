package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRankToNumber(t *testing.T) {
	for _, test := range []struct {
		Rank     *apb.Rank
		Expected uint32
	}{
		{
			Rank: &apb.Rank{
				Tier:     0x10,
				Division: 0x10,
			},
			Expected: 0x00100010,
		},
		{
			Rank: &apb.Rank{
				Tier:     0x1000,
				Division: 0x1000,
			},
			Expected: 0x10001000,
		},
	} {
		assert.Equal(t, test.Expected, RankToNumber(test.Rank))
	}
}

func TestRankFromNumber(t *testing.T) {
	for _, test := range []struct {
		Number   uint32
		Expected Rank
	}{
		{0x00100010, Rank{0x10, 0x10}},
		{0x9010d010, Rank{0x9010, 0xd010}},
	} {
		assert.Equal(t, test.Expected, RankFromNumber(test.Number))
	}
}

func TestParseRank(t *testing.T) {
	for _, test := range []struct {
		Tier     string
		Division string
		Out      uint32
	}{
		{TierBronze, DivisionIII, 0x00100030},
		{TierChallenger, DivisionI, 0x00700050},
	} {
		r, err := ParseRank(test.Tier, test.Division)
		assert.Nil(t, err)
		assert.Equal(t, *r, RankFromNumber(test.Out))
	}
}
