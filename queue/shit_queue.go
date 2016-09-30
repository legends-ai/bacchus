package queue

// ShitQueue is a shitty queue that doesn't do what it's supposed to.
type ShitQueue struct {
	c chan interface{}
}

// NewShitQueue creates a new ShitQueue.
func NewShitQueue() Queue {
	return &ShitQueue{
		c: make(chan interface{}, 1E7),
	}
}

func (q *ShitQueue) Add(in interface{}, features interface{}) {
	q.c <- in
}

func (q *ShitQueue) Poll() interface{} {
	return <-q.c
}
