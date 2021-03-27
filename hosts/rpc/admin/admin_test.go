package admin_test

import (
	"context"
	"testing"

	rpc_admin "github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/admin"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestAdminServer_Config(t *testing.T) {
	t.Parallel()
	srv, _ := tsrv.RPCHost(t, []rpc.RegisterServer{
		admin.RegisterServer,
	})
	client := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return rpc_admin.NewAdminClient(cc)
	}).(rpc_admin.AdminClient)
	ctx := context.Background()
	res, err := client.Settings(ctx, &rpc_admin.SettingsRequest{})
	assert.NoError(t, err)
	test := srv.HealthCheck()
	assert.Equal(t, test.Name, res.Name)
	assert.Equal(t, test.Version, res.Version)
	assert.Equal(t, test.Status, res.Status)
}

func TestAdminServer_Codes(t *testing.T) {
	t.Parallel()
	srv, _ := tsrv.RPCHost(t, []rpc.RegisterServer{
		admin.RegisterServer,
	})
	client := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return rpc_admin.NewAdminClient(cc)
	}).(rpc_admin.AdminClient)
	ctx := context.Background()
	list, err := client.CreateSignupCodes(ctx, &rpc_admin.CreateSignupCodesRequest{
		Uses:  1,
		Count: 1,
	})
	assert.NoError(t, err)
	require.NotNil(t, list)
	assert.Len(t, list.GetCodes(), 1)
	code := list.GetCodes()[0]
	sc, err := client.CheckSignupCode(ctx, &rpc_admin.CheckSignupCodeRequest{
		Code: code,
	})
	assert.NoError(t, err)
	assert.True(t, sc.Valid)
	assert.Equal(t, code, sc.Code)
	_, err = client.DeleteSignupCode(ctx, &rpc_admin.DeleteSignupCodeRequest{
		Code: code,
	})
	assert.NoError(t, err)
	_, err = client.CheckSignupCode(ctx, &rpc_admin.CheckSignupCodeRequest{
		Code: code,
	})
	assert.Error(t, err)
}
