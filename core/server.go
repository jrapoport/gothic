package core

import (
	"github.com/jrapoport/gothic/config"
	"github.com/sirupsen/logrus"
)

// Server holds a core api server
type Server struct {
	*API
	logrus.FieldLogger
}

// NewServer returns a new server.
func NewServer(a *API, name string) *Server {
	log := a.log.WithField("server", name)
	return &Server{a, log}
}

// Clone clones the server.
func (s *Server) Clone() *Server {
	return &Server{
		API:         s.API,
		FieldLogger: s.Log(),
	}
}

// Config returns the config for an api server.
func (s *Server) Config() *config.Config {
	return s.config
}

// Log returns the log
func (s *Server) Log() logrus.FieldLogger {
	return s.FieldLogger
}
