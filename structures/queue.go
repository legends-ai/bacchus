package structures

import "fmt"

const (
	maxPending = 1000000
)

// Queue stores visited and unvisited elements backed by a buffered channel
type Queue struct {
	Visited   *Set
	Unvisited *Set
	Channel   chan RegionedString
}

// NewQueue creates a new queue
func NewQueue() *Queue {
	return &Queue{
		Visited:   NewSet(),
		Unvisited: NewSet(),
		Channel:   make(chan RegionedString, maxPending),
	}
}

// Has the element
func (q *Queue) Has(s RegionedString) bool {
	return q.Visited.Has(s) || q.Unvisited.Has(s)
}

// Offer an element if the queue can accept it
func (q *Queue) Offer(s RegionedString) {
	if !q.Has(s) {
		q.ForceOffer(s)
	}
}

// ForceOffer queues an element for visitation
func (q *Queue) ForceOffer(s RegionedString) {
	q.Unvisited.Add(s)
	q.Channel <- s
}

// Start removes an element from the unvisited queue
func (q *Queue) Start(s RegionedString) {
	q.Unvisited.Remove(s)
}

// Complete marks an element as visited
func (q *Queue) Complete(s RegionedString) {
	q.Visited.Add(s)
}

// Print prints the queue contents
func (q *Queue) Print() {
	fmt.Printf("Visited: %v\n", q.Visited.Values())
	fmt.Printf("Unvisited: %v\n", q.Unvisited.Values())
}
