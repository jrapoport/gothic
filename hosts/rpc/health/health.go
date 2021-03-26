package health

import (
	"context"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/protobuf/grpc/rpc/health"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type healthServer struct {
	health.UnimplementedHealthServer
	*rpc.Server
}

var _ health.HealthServer = (*healthServer)(nil)

func newHealthServer(srv *rpc.Server) *healthServer {
	hs := &healthServer{Server: srv}
	hs.FieldLogger = srv.WithField("module", "health")
	return hs
}

// RegisterServer registers a new health server.
func RegisterServer(s *grpc.Server, e *rpc.Server) {
	health.RegisterHealthServer(s, newHealthServer(e))
}

func (s *healthServer) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	return ctx, nil
}

// HealthCheck performs a health check.
func (s *healthServer) HealthCheck(_ context.Context, _ *emptypb.Empty) (*health.HealthCheckResponse, error) {
	hc := s.API.HealthCheck()
	return &health.HealthCheckResponse{
		Name:    hc.Name,
		Version: hc.Version,
		Status:  hc.Status,
	}, nil
}
