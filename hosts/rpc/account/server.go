package account

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jrapoport/gothic/api/grpc/rpc/account"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
)

type server struct {
	account.UnimplementedAccountServer
	*rpc.Server
}

var (
	_ account.AccountServer             = (*server)(nil)
	_ grpc_auth.ServiceAuthFuncOverride = (*server)(nil)
)

func newServer(srv *rpc.Server) *server {
	srv.Logger = srv.WithName("account")
	return &server{Server: srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	account.RegisterAccountServer(s, newServer(srv))
}

func (s *server) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	// we purposely ignore the error here so we'll parse a token
	// if passed in on the call, but not require it to be there.
	ctx, _ = rpc.Authenticate(ctx, s.Config().JWT)
	return ctx, nil
}
