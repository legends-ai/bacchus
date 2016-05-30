package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	s := []string{
		"a", "b", "c", "d", "e",
	}
	r := Chunk(s, 2)
	assert.Equal(t, 3, len(r))
	assert.Equal(t, []string{"a", "b"}, r[0])
	assert.Equal(t, []string{"c", "d"}, r[1])
	assert.Equal(t, []string{"e"}, r[2])
}
