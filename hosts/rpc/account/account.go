package account

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
)

//go:generate protoc -I=. -I=.. --go_out=plugins=grpc:. --go_opt=paths=source_relative account.proto

type accountServer struct {
	*rpc.Server
}

var (
	_ AccountServer                     = (*accountServer)(nil)
	_ grpc_auth.ServiceAuthFuncOverride = (*accountServer)(nil)
)

func newAccountServer(srv *rpc.Server) *accountServer {
	srv.FieldLogger = srv.WithField("module", "account")
	return &accountServer{srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	RegisterAccountServer(s, newAccountServer(srv))
}

func (s *accountServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}
