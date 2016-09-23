package processor

import (
	"time"

	"github.com/Sirupsen/logrus"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
)

// Metrics records processed summoners and matches and logs progress.
type Metrics struct {
	Logger *logrus.Logger `inject:"t"`
	// SummonerRate is the number of processed summoners to log a message.
	SummonerRate int
	// MatchRate is the number of processed matches to log a message.
	MatchRate int

	sn int
	mn int

	sc chan *apb.SummonerId
	mc chan *apb.MatchId
}

// Start starts the metrics.
func (m *Metrics) Start() {
	m.sc = make(chan *apb.SummonerId)
	m.mc = make(chan *apb.MatchId)

	// Show rate
	go func() {
		for range time.Tick(5 * time.Second) {
			sRate := float64(m.sn) / 5.0
			mRate := float64(m.mn) / 5.0
			m.Logger.Infof("Processed %d summoners (%.2f/sec), %d matches (%.2f/sec)", m.sn, sRate, m.mn, mRate)
			m.sn = 0
			m.mn = 0
		}
	}()

	// Process channels
	for {
		select {
		case <-m.sc:
			m.sn += 1
			break
		case <-m.mc:
			m.mn += 1
			break
		}
	}
}

// RecordSummoner records a summoner.
func (m *Metrics) RecordSummoner(id *apb.SummonerId) {
	m.sc <- id
}

// RecordMatch records a match.
func (m *Metrics) RecordMatch(id *apb.MatchId) {
	m.mc <- id
}
