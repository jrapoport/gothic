package rest

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/improbable-eng/go-httpwares/logging/logrus"
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
	s := core.NewHost(a, name, address)
	s.FieldLogger = s.Log().WithField("protocol", "http")
	rt := NewRouter(s.Config())
	rt.UseLogger(s.FieldLogger)
	server := &http.Server{
		Handler: rt,
	}
	if el, ok := s.FieldLogger.(*logrus.Entry); ok {
		server.ErrorLog = http_logrus.AsHttpLogger(el)
	}
	for _, r := range reg {
		srv := NewServer(s.Server.Clone())
		r(server, srv)
	}
	return &Host{s, server}
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
