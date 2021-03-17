package health

//go:generate protoc -I=. --go_out=plugins=grpc:. --go_opt=paths=source_relative health.proto

import (
	"context"

	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type healthServer struct {
	*rpc.Server
}

var _ HealthServer = (*healthServer)(nil)

func newHealthServer(srv *rpc.Server) *healthServer {
	hs := &healthServer{srv}
	hs.FieldLogger = srv.WithField("module", "health")
	return hs
}

// RegisterServer registers a new health server.
func RegisterServer(s *grpc.Server, e *rpc.Server) {
	RegisterHealthServer(s, newHealthServer(e))
}

func (s *healthServer) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	return ctx, nil
}

// HealthCheck performs a health check.
func (s *healthServer) HealthCheck(_ context.Context, _ *emptypb.Empty) (*HealthCheckResponse, error) {
	hc := s.API.HealthCheck()
	return &HealthCheckResponse{
		Name:    hc.Name,
		Version: hc.Version,
		Status:  hc.Status,
	}, nil
}
