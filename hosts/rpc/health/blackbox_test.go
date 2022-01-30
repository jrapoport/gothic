package health_test

import (
	"context"
	"testing"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/health"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthServer_HealthCheck(t *testing.T) {
	t.Parallel()
	srv, _ := tsrv.RPCHost(t, []rpc.RegisterServer{
		health.RegisterServer,
	})
	hc := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return healthpb.NewHealthClient(cc)
	}).(healthpb.HealthClient)
	ctx := context.Background()
	res, err := hc.Check(ctx, &healthpb.HealthCheckRequest{})
	assert.NoError(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, res.GetStatus())
}
