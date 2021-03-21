package admin

import (
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/admin/codes"
	"github.com/jrapoport/gothic/hosts/rpc/admin/settings"
	"google.golang.org/grpc"
)

type adminServer struct {
	*rpc.Server
}

func newAdminServer(srv *rpc.Server) *adminServer {
	srv.FieldLogger = srv.WithField("module", "admin")
	return &adminServer{srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	as := newAdminServer(srv).Server
	settings.RegisterServer(s, as)
	codes.RegisterServer(s, as)
}
