package admin

import (
	"context"
	"testing"

	rpc_admin "github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/config"
	core_ctx "github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// rootContext test rpc root context
func rootContext(c *config.Config) context.Context {
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, c.RootPassword))
	return core_ctx.WithContext(ctx)
}

func TestAdminServer_Config(t *testing.T) {
	t.Parallel()
	srv, _ := tsrv.RPCHost(t, []rpc.RegisterServer{
		RegisterServer,
	})
	client := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return rpc_admin.NewAdminClient(cc)
	}).(rpc_admin.AdminClient)
	ctx := metadata.NewOutgoingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, srv.Config().RootPassword))
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
		RegisterServer,
	}, rpc.Authentication())
	client := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return rpc_admin.NewAdminClient(cc)
	}).(rpc_admin.AdminClient)
	_, tok := tcore.TestUser(t, srv.API, "", true)
	ctx := tsrv.RPCAuthContext(t, srv.Config(), tok)
	// authenticated call (success)
	req := &rpc_admin.CreateSignupCodesRequest{
		Uses:  1,
		Count: 1,
	}
	list, err := client.CreateSignupCodes(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, list)
	assert.Len(t, list.GetCodes(), 1)
	code := list.GetCodes()[0]
	creq := &rpc_admin.CheckSignupCodeRequest{
		Code: code,
	}
	sc, err := client.CheckSignupCode(ctx, creq)
	assert.NoError(t, err)
	assert.True(t, sc.Valid)
	assert.Equal(t, code, sc.Code)
	dreq := &rpc_admin.DeleteSignupCodeRequest{
		Code: code,
	}
	_, err = client.DeleteSignupCode(ctx, dreq)
	assert.NoError(t, err)
	creq = &rpc_admin.CheckSignupCodeRequest{
		Code: code,
	}
	_, err = client.CheckSignupCode(ctx, creq)
	assert.Error(t, err)
}
