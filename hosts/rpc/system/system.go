package system

import (
	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
)

type systemServer struct {
	system.UnimplementedSystemServer
	*rpc.Server
}

var _ system.SystemServer = (*systemServer)(nil)

func newSystemServer(srv *rpc.Server) *systemServer {
	srv.Logger = srv.WithName("user")
	return &systemServer{Server: srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	system.RegisterSystemServer(s, newSystemServer(srv))
}
