package processor

import (
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
	for {
		select {
		case id := <-m.sc:
			m.sn += 1
			if m.sn%m.SummonerRate == 0 {
				m.Logger.Infof("Processed %d summoners (%s)", m.sn, id.String())
			}
			break
		case id := <-m.mc:
			m.mn += 1
			if m.mn%m.MatchRate == 0 {
				m.Logger.Infof("Processed %d matches (%s)", m.mn, id.String())
			}
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
