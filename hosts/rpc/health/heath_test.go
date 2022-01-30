package health

import (
	"context"
	"testing"

	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthServer_HealthCheck(t *testing.T) {
	t.Parallel()
	s, _ := tsrv.RPCServer(t, false)
	srv := newHealthServer(s)
	ctx := context.Background()
	res, err := srv.Check(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, res.GetStatus())
}
