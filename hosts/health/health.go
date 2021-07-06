package health

import (
	"net/http"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/health"
	"github.com/jrapoport/gothic/hosts/rest"
)

// Name is the name of the health host.
const Name = "health"

// Endpoint is the endpoint for the health host.
const Endpoint = "/*"

// Health check response
type Health struct {
	health.Health
	Hosts map[string]Status `json:"hosts"`
}

// Status is returned by the health host.
type Status struct {
	Name    string `json:"name"`
	Address string `json:"address,omitempty"`
	Online  bool   `json:"online"`
}

func hostStatus(h core.Hosted) Status {
	if h == nil {
		return Status{}
	}
	return Status{
		Name:    h.Name(),
		Address: h.Address(),
		Online:  h.Online(),
	}
}

// NewHealthHost creates a new http host.
func NewHealthHost(a *core.API, addr string, hosted map[string]core.Hosted) core.Hosted {
	s := rest.NewHost(a, Name, addr,
		[]rest.RegisterServer{
			func(s *http.Server, server *rest.Server) {
				registerServer(s, server, hosted)
			},
		})
	return s
}

func registerServer(s *http.Server, srv *rest.Server, hosted map[string]core.Hosted) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		hc := srv.HealthCheck()
		status := map[string]Status{}
		for k, v := range hosted {
			if k == Name {
				continue
			}
			stat := hostStatus(v)
			if stat.Name == "" {
				stat.Name = k
			}
			status[k] = stat
		}
		res := &Health{
			Health: hc,
			Hosts:  status,
		}
		srv.Response(w, res)
	}
	if r, ok := s.Handler.(*rest.Router); ok {
		r.Get(Endpoint, handler)
	}
}
