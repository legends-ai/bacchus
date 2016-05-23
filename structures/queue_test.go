package structures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestQueuePollOffer tests offering and polling a queue.
func TestQueuePollOffer(t *testing.T) {
	s := &QueueSettings{
		Concurrency: 100,
	}
	q := s.Create()
	q.Offer("a")
	q.Offer("b")
	assert.True(t, q.Has("a"))
	assert.True(t, q.Has("b"))
	assert.False(t, q.Has("c"))
	assert.Equal(t, q.Poll(), "a")
	assert.Equal(t, q.Poll(), "b")
}

func TestQueueComplete(t *testing.T) {
	s := &QueueSettings{
		Concurrency: 100,
	}
	q := s.Create()
	q.Complete("test")
	assert.True(t, q.Has("test"))
}

func TestQueueOfferBlockedIfComplete(t *testing.T) {
	s := &QueueSettings{
		Concurrency: 100,
	}
	q := s.Create()
	q.Complete("test")
	q.Offer("test")
	assert.False(t, q.Unvisited.Has("test"))
}
