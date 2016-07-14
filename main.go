package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/bacchus/lib"
	"github.com/simplyianm/bacchus/models"
	"github.com/simplyianm/bacchus/processor"
)

const (
	concurrency = 10
)

func main() {
	inject := lib.NewInjector()
	inject.Invoke(startProcessors)
	select {}
}

func startProcessors(s *processor.Summoners, m *processor.Matches, logger *logrus.Logger) {
	for i := 0; i < concurrency; i++ {
		logger.Infof("Starting summoner processor %d", i)
		go s.Start()
	}
	for i := 0; i < concurrency; i++ {
		logger.Infof("Starting match processor %d", i)
		go m.Start()
	}
	// Offer Aditi
	s.Offer(models.SummonerID{"na", 32875076})
}
