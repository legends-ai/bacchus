package queue

// Queue is a priority queue.
type Queue interface {
	// Add adds an element to the queue.
	// Features is attributes of the queue. It can be whatever you want
	// it to be. You just have to believe!
	Add(in interface{}, features interface{})

	// Poll gets the next element of the queue to process.
	Poll() interface{}
}
