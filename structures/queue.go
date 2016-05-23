package structures

import "fmt"

// QueueSettings represents settings to build a queue
type QueueSettings struct {
	Concurrency int
}

// Queue stores visited and unvisited elements backed by a buffered channel
type Queue struct {
	QueueSettings
	Visited   StringSet
	Unvisited StringSet
	Channel   chan string
}

// Create a queue
func (q QueueSettings) Create() *Queue {
	return &Queue{
		QueueSettings: q,
		Visited:       StringSet{},
		Unvisited:     StringSet{},
		Channel:       make(chan string, q.Concurrency),
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

// Poll polls the queue for the next element
func (q *Queue) Poll() string {
	ret := <-q.Channel
	q.Unvisited.Remove(ret)
	return ret
}

// ForceOffer queues an element for visitation
func (q *Queue) ForceOffer(s string) {
	q.Unvisited.Add(s)
	q.Channel <- s
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
