package account_test

import (
	"context"
	"testing"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/account"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestLogin(t *testing.T) {
	srv, _ := tsrv.RPCHost(t, []rpc.RegisterServer{
		account.RegisterServer,
	})
	t.Log(srv.Address())
	client := tsrv.RPCClient(t, srv.Address(), func(cc grpc.ClientConnInterface) interface{} {
		return account.NewAccountClient(cc)
	}).(account.AccountClient)
	// new request context
	ctx := context.Background()
	// add key-value pairs of metadata to context
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("key1", "val1", "key2", "val2"),
	)
	req := &account.LoginRequest{}
	_, err := client.Login(ctx, req)
	assert.Error(t, err)
}
