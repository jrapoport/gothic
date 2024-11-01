package admin

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/admin/audit"
	"github.com/jrapoport/gothic/hosts/rest/admin/codes"
	"github.com/jrapoport/gothic/hosts/rest/admin/settings"
	"github.com/jrapoport/gothic/hosts/rest/admin/users"
	"github.com/jrapoport/gothic/hosts/rest/modules/invite"
)

// Admin is the admin endpoint.
const Admin = "/admin"

type adminServer struct {
	*rest.Server
}

func newAdminServer(srv *rest.Server) *adminServer {
	srv.Logger = srv.WithName("admin")
	return &adminServer{srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newAdminServer(srv))
}

func register(s *http.Server, srv *adminServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *adminServer) addRoutes(r *rest.Router) {
	r.Authenticated().Admin().Route(Admin, func(rt *rest.Router) {
		audit.RegisterServer(&http.Server{Handler: rt}, s.Clone())
		invite.RegisterServer(&http.Server{Handler: rt}, s.Clone())
		settings.RegisterServer(&http.Server{Handler: rt}, s.Clone())
		codes.RegisterServer(&http.Server{Handler: rt}, s.Clone())
		users.RegisterServer(&http.Server{Handler: rt}, s.Clone())
	})
}
