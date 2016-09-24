package server

import (
	apb "github.com/asunaio/bacchus/gen-go/asuna"
	"golang.org/x/net/context"
)

type Server struct {
}

// Ingest implements the Bacchus.Ingest RPC endpoint.
func (s *Server) Ingest(
	ctx context.Context, in *apb.IngestRequest,
) (*apb.IngestResponse, error) {
	// TODO(igm): implement
	return nil, nil
}
