package spider

import (
	"bytes"
	"sync"

	"github.com/Sirupsen/logrus"
)

const maxBatchSizeBytes = 50 * 1024 * 1024 // 50mib

// Batcher aggregates input JSON lines and writes them to S3.
type Batcher struct {
	Logger *logrus.Logger `inject:"t"`
	mutex  sync.Mutex
	batch  *bytes.Buffer
	input  chan string
}

// NewBatchr constructs a batcher
func NewBatcher() *Batcher {
	return &Batcher{
		batch: &bytes.Buffer{},
		input: make(chan string),
	}
}

// Add adds another line to the batcher
func (b *Batcher) Add(json string) {
	b.input <- json
}

// Start starts the batcher
func (b *Batcher) Start() {
	for {
		json := <-b.input
		b.mutex.Lock()
		b.batch.Write([]byte(json + "\n"))
		b.mutex.Unlock()
		if b.batch.Len() > maxBatchSizeBytes {
			go b.nextBatch()
		}
	}
}

// nextBatch creates a new batch
func (b *Batcher) nextBatch() {
	// Move the buffer reference into the function so we don't block other calls
	b.mutex.Lock()
	buf := b.batch
	b.batch = &bytes.Buffer{}
	b.mutex.Unlock()
	if err := b.writeToS3(buf); err != nil {
		b.Logger.Errorf("Could not write batch: %v", err)
	}
}

// writeToS3 writes stuff to S3
func (b *Batcher) writeToS3(buf *bytes.Buffer) error {
	// TODO(simplyianm): implementation
	return nil
}
