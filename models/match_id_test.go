package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchIDString(t *testing.T) {
	assert.Equal(t, MatchID{
		Region: "na",
		ID:     1738,
	}.String(), "na/1738")
}
