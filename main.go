package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"google.golang.org/grpc"

	"github.com/Sirupsen/logrus"
	"github.com/asunaio/bacchus/config"
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/lib"
	"github.com/asunaio/bacchus/processor"
	"github.com/asunaio/bacchus/server"
	"github.com/simplyianm/inject"
)

func main() {
	inject := lib.NewInjector()

	_, err := inject.Invoke(startProcessors)
	if err != nil {
		log.Fatalf("Could not start processors: %v", err)
	}

	_, err = inject.Invoke(startServer)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func startProcessors(
	cfg *config.AppConfig, s *processor.Summoners, m *processor.Matches, logger *logrus.Logger,
) {
	go func() {
		for i := 0; i < cfg.Concurrency; i++ {
			go s.Start()
		}
		for i := 0; i < cfg.Concurrency; i++ {
			go m.Start()
		}
		s.Seed()
	}()
}

func startServer(logger *logrus.Logger, cfg *config.AppConfig, injector inject.Injector) {
	// Listen on port
	port := fmt.Sprintf(":%d", cfg.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Setup gRPC server
	s := grpc.NewServer()
	serv := &server.Server{}

	_, err = injector.ApplyMap(serv)
	if err != nil {
		logger.Fatalf("Could not inject server: %v", err)
	}

	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		mp := fmt.Sprintf(":%d", cfg.MonitorPort)
		logger.Infof("Monitor listening on %s", mp)
		http.ListenAndServe(mp, nil)
	}()

	apb.RegisterBacchusServer(s, serv)
	logger.Infof("Listening on %s", port)
	s.Serve(lis)
}
