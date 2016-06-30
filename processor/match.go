package processor

// MatchID identifies a match.
type MatchID struct {
	Region string
	ID     int
}

// Matches is the processor for matches.
type Matches struct {
	c chan MatchID
}

// NewMatches creates a new processor.Matches.
func NewMatches() *Matches {
	return &Matches{
		c: make(chan MatchID),
	}
}

// Offer offers a match to the queue which may accept it.
func (m *Matches) Offer(id MatchID) {
	// if key exists in cassandra return
}

// Start starts processing matches.
func (m *Matches) Start() {
	for {
		id, ok := <-m.c
		if !ok {
			return
		}
		m.process(id)
	}
}

func (m *Matches) process(id MatchID) {
}
