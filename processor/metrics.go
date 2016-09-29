package processor

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

// Metrics records processed summoners and matches and logs progress.
type Metrics struct {
	Logger *logrus.Logger `inject:"t"`

	reqCt   map[string]int
	reqs    chan string
	reqCtMu sync.Mutex
}

// Start starts the metrics.
func (m *Metrics) Start() {
	m.reqCt = map[string]int{}
	m.reqs = make(chan string)

	// Show rate
	go func() {
		for range time.Tick(5 * time.Second) {
			m.Logger.Infof("===")
			m.reqCtMu.Lock()
			total := 0
			for reqType, ct := range m.reqCt {
				m.Logger.Infof("- %s: %d (%.2f/sec)", reqType, ct, float64(ct)/5.0)
				total += ct
			}
			m.Logger.Infof("TOTAL: %d (%.2f/sec)", total, float64(total)/5.0)
			m.reqCt = map[string]int{}
			m.reqCtMu.Unlock()
		}
	}()

	// Process channels
	for endpoint := range m.reqs {
		m.reqCtMu.Lock()
		m.reqCt[endpoint] += 1
		m.reqCtMu.Unlock()
	}
}

// RecordSummoner records a summoner.
func (m *Metrics) Record(endpoint string) {
	m.reqs <- endpoint
}
