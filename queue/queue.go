package queue

// Queue is a priority queue.
type Queue interface {
	// Add adds an element to the queue.
	// Context is an arbitrary value which may influence the priority of the element within the queue.
	Add(in interface{}, context interface{})

	// Poll gets the next element of the queue to process.
	Poll() interface{}
}
