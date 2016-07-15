package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRankToNumber(t *testing.T) {
	for _, test := range []struct {
		Tier     uint16
		Division uint16
		Expected uint32
	}{
		{0x10, 0x10, 0x00100010},
		{0x1000, 0x1009, 0x10001009},
	} {
		assert.Equal(t, test.Expected, Rank{test.Tier, test.Division}.ToNumber())
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