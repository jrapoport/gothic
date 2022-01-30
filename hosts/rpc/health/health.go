package health

import (
	"context"

	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type healthServer struct {
	healthpb.UnimplementedHealthServer
	*rpc.Server
}

var _ healthpb.HealthServer = (*healthServer)(nil)

func newHealthServer(srv *rpc.Server) *healthServer {
	hs := &healthServer{Server: srv}
	hs.Logger = srv.WithName("health")
	return hs
}

// RegisterServer registers a new health server.
func RegisterServer(s *grpc.Server, e *rpc.Server) {
	healthpb.RegisterHealthServer(s, newHealthServer(e))
}

func (s *healthServer) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	return ctx, nil
}

// Check performs a health check.
func (s *healthServer) Check(_ context.Context, _ *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is ignored for now
func (s *healthServer) Watch(_ *healthpb.HealthCheckRequest, _ healthpb.Health_WatchServer) error {
	//TODO implement me
	panic("not implemented")
}
