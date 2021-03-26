package account

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jrapoport/gothic/api/grpc/rpc/account"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
)

type accountServer struct {
	account.UnimplementedAccountServer
	*rpc.Server
}

var (
	_ account.AccountServer             = (*accountServer)(nil)
	_ grpc_auth.ServiceAuthFuncOverride = (*accountServer)(nil)
)

func newAccountServer(srv *rpc.Server) *accountServer {
	srv.FieldLogger = srv.WithField("module", "account")
	return &accountServer{Server: srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	account.RegisterAccountServer(s, newAccountServer(srv))
}

func (s *accountServer) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	// we purposely ignore the error here so we'll parse a token
	// if passed in on the call, but not require it to be there.
	ctx, _ = rpc.Authenticate(ctx, s.Config().JWT)
	return ctx, nil
}
