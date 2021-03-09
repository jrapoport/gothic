package account

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account/auth"
	"github.com/jrapoport/gothic/hosts/rest/account/confirm"
	"github.com/jrapoport/gothic/hosts/rest/account/login"
	"github.com/jrapoport/gothic/hosts/rest/account/password"
	"github.com/jrapoport/gothic/hosts/rest/account/signup"
)

// Endpoint is the account endpoint
const Endpoint = "/account"

type accountServer struct {
	*rest.Server
}

// newAccountServer returns a new account server.
func newAccountServer(srv *rest.Server) *accountServer {
	srv.FieldLogger = srv.WithField("service", "account")
	return &accountServer{srv}
}

// RegisterServer registers an account server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newAccountServer(srv))
}

func register(s *http.Server, srv *accountServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *accountServer) addRoutes(r *rest.Router) {
	r.RateLimit().Route(Endpoint, func(rt *rest.Router) {
		auth.RegisterServer(&http.Server{Handler: rt}, s.Clone())
		confirm.RegisterServer(&http.Server{Handler: rt}, s.Clone())
		login.RegisterServer(&http.Server{Handler: rt}, s.Clone())
		password.RegisterServer(&http.Server{Handler: rt}, s.Clone())
		signup.RegisterServer(&http.Server{Handler: rt}, s.Clone())
	})
}
