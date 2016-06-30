package processor

// SummonerID identifies a summoner.
type SummonerID struct {
	Region     string
	SummonerID int
}

// Queues is the processor for queues.
type Summoners struct {
	c      chan SummonerID
	exists map[SummonerID]bool
}

// NewSummoners creates a new processor.Summoners.
func NewSummoners() *Summoners {
	return &Summoners{
		c:      make(chan SummonerID),
		exists: map[SummonerID]bool{},
	}
}

// Offer offers a summoner to the queue which may accept it.
func (q *Summoners) Offer(s SummonerID) {
	if q.exists[s] {
		return
	}
	q.c <- s
}

// Start starts processing summoners.
func (q *Summoners) Start() {
	for {
		s, ok := <-q.c
		if !ok {
			return
		}
		q.process(s)
	}
}

func (q *Summoners) process(m SummonerID) {
}
