package account

import (
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
)

//go:generate protoc -I=. -I=../.. --go_out=plugins=grpc:. --go_opt=paths=source_relative account.proto

type accountServer struct {
	*rpc.Server
}

var _ AccountServer = (*accountServer)(nil)

// NewAccountServer returns a new account server
func newAccountServer(srv *rpc.Server) *accountServer {
	srv.FieldLogger = srv.WithField("module", "account")
	return &accountServer{srv}
}

// RegisterServer registers a new account server.
func RegisterServer(s *grpc.Server, e *rpc.Server) {
	RegisterAccountServer(s, newAccountServer(e))
}

func (s *accountServer) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	return ctx, nil
}
