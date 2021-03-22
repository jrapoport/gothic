package settings

import (
	"net/http"

	"github.com/jrapoport/gothic/hosts/rest"
)

// Settings is the config endpoint
const Settings = "/settings"

type settingsServer struct {
	*rest.Server
}

// NewSettingsServer returns a new config rest server.
func newSettingsServer(srv *rest.Server) *settingsServer {
	srv.FieldLogger = srv.WithField("module", "settings")
	return &settingsServer{srv}
}

// RegisterServer registers a new config server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newSettingsServer(srv))
}

func register(s *http.Server, srv *settingsServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *settingsServer) addRoutes(r *rest.Router) {
	r.Get(Settings, s.Settings)
}

func (s *settingsServer) Settings(w http.ResponseWriter, r *http.Request) {
	s.Response(w, s.API.Settings())
}
