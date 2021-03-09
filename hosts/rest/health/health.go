package health

import (
	"net/http"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/hosts/rest"
)

// Endpoint is the health check endpoint.
const Endpoint = config.HealthEndpoint

type healthServer struct {
	*rest.Server
}

func newHealthServer(srv *rest.Server) *healthServer {
	srv.FieldLogger = srv.WithField("module", "health")
	return &healthServer{srv}
}

// RegisterServer registers a new health server.
func RegisterServer(s *http.Server, srv *rest.Server) {
	register(s, newHealthServer(srv))
}

func register(s *http.Server, srv *healthServer) {
	if r, ok := s.Handler.(*rest.Router); ok {
		srv.addRoutes(r)
	}
}

func (s *healthServer) addRoutes(r *rest.Router) {
	r.Get(Endpoint, s.HealthCheck)
}

// HealthCheck performs a health check.
func (s *healthServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	s.Response(w, s.API.HealthCheck())
}
