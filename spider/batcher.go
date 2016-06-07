package spider

import "bytes"

type Batcher struct {
	batch *bytes.Buffer
	input chan string
}

func NewBatcher() *Batcher {
	return &Batcher{
		batch: &bytes.Buffer{},
		input: make(chan string),
	}
}

func (b *Batcher) Add(json string) {
}

// Start starts the batcher
func (b *Batcher) Start() {
	for {
		json, more := <-b.input
		b.batch.Write([]byte(json + "\n"))
	}
}
