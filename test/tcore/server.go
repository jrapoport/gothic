package tcore

import (
	"testing"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/test/tconf"
)

// Server for testing.
func Server(t *testing.T, smtp bool) (*core.Server, *tconf.SMTPMock) {
	a, _, mock := API(t, smtp)
	s := core.NewServer(a, "test")
	return s, mock
}
