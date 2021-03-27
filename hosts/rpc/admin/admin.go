package admin

import (
	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
)

type adminServer struct {
	admin.UnimplementedAdminServer
	*rpc.Server
}

func newAdminServer(srv *rpc.Server) *adminServer {
	srv.FieldLogger = srv.WithField("module", "admin")
	return &adminServer{Server: srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	admin.RegisterAdminServer(s, newAdminServer(srv))
}
