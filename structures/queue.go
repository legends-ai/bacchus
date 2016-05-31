package structures

import "fmt"

const (
	maxPending = 1000000
)

// Queue stores visited and unvisited elements backed by a buffered channel
type Queue struct {
	Visited   *StringSet
	Unvisited *StringSet
	Channel   chan string
}

// NewQueue creates a new queue
func NewQueue() *Queue {
	return &Queue{
		Visited:   NewStringSet(),
		Unvisited: NewStringSet(),
		Channel:   make(chan string, maxPending),
	}
}

// Has the element
func (q *Queue) Has(s string) bool {
	return q.Visited.Has(s) || q.Unvisited.Has(s)
}

// Offer an element if the queue can accept it
func (q *Queue) Offer(s string) {
	if !q.Has(s) {
		q.ForceOffer(s)
	}
}

// ForceOffer queues an element for visitation
func (q *Queue) ForceOffer(s string) {
	q.Unvisited.Add(s)
	q.Channel <- s
}

// Start removes an element from the unvisited queue
func (q *Queue) Start(s string) {
	q.Unvisited.Remove(s)
}

// Complete marks an element as visited
func (q *Queue) Complete(s string) {
	q.Visited.Add(s)
}

// Print prints the queue contents
func (q *Queue) Print() {
	fmt.Printf("Visited: %v\n", q.Visited.Values())
	fmt.Printf("Unvisited: %v\n", q.Unvisited.Values())
}
