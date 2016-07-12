package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummonerIDString(t *testing.T) {
	assert.Equal(t, SummonerID{
		Region: "na",
		ID:     1738,
	}.String(), "na/1738")
}
