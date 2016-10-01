package server

import (
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"github.com/asunaio/bacchus/processor"
	"golang.org/x/net/context"
)

type Server struct {
	Matches   *processor.Matches   `inject:"t"`
	Summoners *processor.Summoners `inject:"t"`
}

// Ingest implements the Bacchus.Ingest RPC endpoint.
func (s *Server) Ingest(
	ctx context.Context, in *apb.IngestRequest,
) (*apb.IngestResponse, error) {
	// TODO(igm): implement
	return nil, nil
}
