package structures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Integration test
func TestStringSet(t *testing.T) {
	set := StringSet{}
	set.Add("a")
	assert.True(t, set.Has("a"))
	assert.False(t, set.Has("b"))
	set.Remove("a")
	assert.False(t, set.Has("a"))
}
