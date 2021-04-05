package rest

import (
	"context"
	"errors"
	"net"
	"net/http"

	http_logrus "github.com/improbable-eng/go-httpwares/logging/logrus"
	"github.com/jrapoport/gothic/core"
	"github.com/sirupsen/logrus"
)

// RegisterServer registers an http server to the host.
type RegisterServer func(s *http.Server, srv *Server)

// Host represents a HTTP host.
type Host struct {
	*core.Host
	server *http.Server
}

var _ core.Hosted = (*Host)(nil)

// NewHost creates a new Host.
func NewHost(a *core.API, name string, address string, reg []RegisterServer) *Host {
	h := core.NewHost(a, name, address)
	h.Logger = h.Log().WithName("protocol-http")
	rt := NewRouter(h.Config())
	//rt.UseLogger(h.Logger)
	server := &http.Server{
		Handler: rt,
	}
	log := httpLogger(h.Config().Level)
	server.ErrorLog = http_logrus.AsHttpLogger(log)
	for _, r := range reg {
		srv := NewServer(h.Server.Clone())
		r(server, srv)
	}
	return &Host{h, server}
}

// ListenAndServe starts the http server.
func (s *Host) ListenAndServe() error {
	s.Start(func(lis net.Listener) error {
		err := s.server.Serve(lis)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	return s.Host.ListenAndServe()
}

// Shutdown stops the http server.
func (s *Host) Shutdown() error {
	s.Stop(func(ctx context.Context) error {
		s.server.SetKeepAlivesEnabled(false)
		return s.server.Shutdown(ctx)
	})
	return s.Host.Shutdown()
}

func httpLogger(level string) *logrus.Entry {
	lvl, _ := logrus.ParseLevel(level)
	l := logrus.New()
	l.SetLevel(lvl)
	return l.WithField("protocol", "grpc")
}
