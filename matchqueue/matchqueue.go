package matchqueue

// MatchQueue represents a queue to process matches.
type MatchQueue struct {
	c chan string
}

// New creates a new MatchQueue.
func New() *MatchQueue {
	return &MatchQueue{
		c: make(chan string),
	}
}

func (q *MatchQueue) Offer(matchId string) {
	// if key exists in cassandra return
}

func (q *MatchQueue) Process(matchId string) {
}
