package config

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
)

// Endpoint is the config endpoint
const Endpoint = "/config"

type configServer struct {
	*rest.Server
}

// NewConfigServer returns a new config rest server.
func newConfigServer(srv *rest.Server) *configServer {
	srv.FieldLogger = srv.WithField("module", "config")
	return &configServer{srv}
}

// RegisterServer registers a new config server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newConfigServer(srv))
}

func register(s *http.Server, srv *configServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *configServer) addRoutes(r *rest.Router) {
	r.Get(Endpoint, s.Config)
}

func (s *configServer) Config(w http.ResponseWriter, r *http.Request) {
	s.Response(w, s.API.Settings())
}
