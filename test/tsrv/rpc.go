package tsrv

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/jwt"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RPCServer rpc server for tests
func RPCServer(t *testing.T, smtp bool) (*rpc.Server, *tconf.SMTPMock) {
	s, mock := tcore.Server(t, smtp)
	return rpc.NewServer(s), mock
}

// RPCHost an rpc server for tests.
func RPCHost(t *testing.T, reg []rpc.RegisterServer, opt ...grpc.ServerOption) (*rpc.Host, *tconf.SMTPMock) {
	a, _, mock := tcore.API(t, true)
	h := rpc.NewHost(a, "test-rpc", "127.0.0.1:0", reg, opt...)
	err := h.ListenAndServe()
	require.NoError(t, err)
	t.Cleanup(func() {
		err = h.Shutdown()
		assert.NoError(t, err)
	})
	return h, mock
}

// RPCClient rpc client for tests
func RPCClient(t *testing.T, addr string, cli func(cc grpc.ClientConnInterface) interface{}) interface{} {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(addr, opts...)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = conn.Close()
		assert.NoError(t, err)
	})
	return cli(conn)
}

// RPCAuthContext test rpc auth context
func RPCAuthContext(t *testing.T, c *config.Config, tok string) context.Context {
	claims, err := jwt.ParseUserClaims(c.JWT, tok)
	require.NoError(t, err)
	require.NotNil(t, claims)
	ctx := metadata.NewOutgoingContext(context.Background(),
		metadata.Pairs(rpc.Authorization, rpc.BearerScheme+" "+tok))
	return context.WithContext(rpc.WithClaims(ctx, claims))
}
