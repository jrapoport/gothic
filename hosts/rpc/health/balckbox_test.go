package health_test

import (
	"context"
	"testing"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/health"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestHealthServer_HealthCheck(t *testing.T) {
	srv, _ := tsrv.RPCHost(t, []rpc.RegisterServer{
		health.RegisterServer,
	})
	hc := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return health.NewHealthClient(cc)
	}).(health.HealthClient)
	ctx := context.Background()
	req := &health.HealthCheckRequest{}
	res, err := hc.HealthCheck(ctx, req)
	assert.NoError(t, err)
	test := srv.HealthCheck()
	assert.Equal(t, test.Name, res.Name)
	assert.Equal(t, test.Version, res.Version)
	assert.Equal(t, test.Status, res.Status)
}
