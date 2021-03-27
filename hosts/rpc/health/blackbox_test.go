package health_test

import (
	"context"
	"testing"

	rpc_health "github.com/jrapoport/gothic/api/grpc/rpc/health"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/health"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestHealthServer_HealthCheck(t *testing.T) {
	t.Parallel()
	srv, _ := tsrv.RPCHost(t, []rpc.RegisterServer{
		health.RegisterServer,
	})
	hc := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return rpc_health.NewHealthClient(cc)
	}).(rpc_health.HealthClient)
	ctx := context.Background()
	res, err := hc.HealthCheck(ctx, &emptypb.Empty{})
	assert.NoError(t, err)
	test := srv.HealthCheck()
	assert.Equal(t, test.Name, res.Name)
	assert.Equal(t, test.Version, res.Version)
	assert.Equal(t, test.Status, res.Status)
}
