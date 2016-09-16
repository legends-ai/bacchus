package processor

import (
	"github.com/Sirupsen/logrus"
	"github.com/asunaio/bacchus/models"
)

// Metrics records processed summoners and matches and logs progress.
type Metrics struct {
	Logger *logrus.Logger `inject:"t"`
	// SummonerRate is the number of processed summoners to log a message.
	SummonerRate int
	// MatchRate is the number of processed matches to log a message.
	MatchRate int
	sn        int
	mn        int
	sc        chan models.SummonerID
	mc        chan models.MatchID
}

// Start starts the metrics.
func (m *Metrics) Start() {
	m.sc = make(chan models.SummonerID)
	m.mc = make(chan models.MatchID)
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
func (m *Metrics) RecordSummoner(id models.SummonerID) {
	m.sc <- id
}

// RecordMatch records a match.
func (m *Metrics) RecordMatch(id models.MatchID) {
	m.mc <- id
}
