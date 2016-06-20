package structures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Integration test
func TestSet(t *testing.T) {
	set := NewSet()
	set.Add("a")
	assert.True(t, set.Has("a"))
	assert.False(t, set.Has("b"))
	set.Remove("a")
	assert.False(t, set.Has("a"))
}
