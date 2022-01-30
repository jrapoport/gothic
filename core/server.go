package core

import (
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/log"
)

// Server holds a core api server
type Server struct {
	*API
	log.Logger
}

// NewServer returns a new server.
func NewServer(a *API, name string) *Server {
	l := a.log.WithName("server-" + name)
	return &Server{a, l}
}

// Clone clones the server.
func (s *Server) Clone() *Server {
	return &Server{
		API:    s.API,
		Logger: s.Log(),
	}
}

// Config returns the config for an api server.
func (s *Server) Config() *config.Config {
	return s.config
}

// Log returns the log
func (s *Server) Log() log.Logger {
	return s.Logger
}
