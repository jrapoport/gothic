package admin

import (
	"context"
	"testing"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/admin/codes"
	"github.com/jrapoport/gothic/hosts/rpc/admin/config"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestAdminServer_Config(t *testing.T) {
	srv, _ := tsrv.RPCHost(t, []rpc.RegisterServer{
		RegisterServer,
	})
	client := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return config.NewConfigClient(cc)
	}).(config.ConfigClient)
	ctx := context.Background()
	res, err := client.Settings(ctx, &config.SettingsRequest{})
	assert.NoError(t, err)
	test := srv.HealthCheck()
	assert.Equal(t, test.Name, res.Name)
	assert.Equal(t, test.Version, res.Version)
	assert.Equal(t, test.Status, res.Status)
}

func TestAdminServer_Codes(t *testing.T) {
	srv, _ := tsrv.RPCHost(t, []rpc.RegisterServer{
		RegisterServer,
	})
	client := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return codes.NewCodesClient(cc)
	}).(codes.CodesClient)
	ctx := context.Background()
	list, err := client.CreateSignupCodes(ctx, &codes.CreateSignupCodesRequest{
		Uses:  1,
		Count: 1,
	})
	assert.NoError(t, err)
	require.NotNil(t, list)
	assert.Len(t, list.GetCodes(), 1)
	code := list.GetCodes()[0]
	sc, err := client.CheckSignupCode(ctx, &codes.CheckSignupCodeRequest{
		Code: code,
	})
	assert.NoError(t, err)
	assert.True(t, sc.Usable)
	assert.Equal(t, code, sc.Code)
	_, err = client.VoidSignupCode(ctx, &codes.VoidSignupCodeRequest{
		Code: code,
	})
	assert.NoError(t, err)
	_, err = client.CheckSignupCode(ctx, &codes.CheckSignupCodeRequest{
		Code: code,
	})
	assert.Error(t, err)
}
