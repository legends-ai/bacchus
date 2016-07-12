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
		{0x1000, 0x1000, 0x10001000},
	} {
		assert.Equal(t, test.Expected, Rank{test.Tier, test.Division}.ToNumber())
	}
}
