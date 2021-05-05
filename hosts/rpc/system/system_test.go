package system

import (
	"testing"

	"github.com/jrapoport/gothic/test/tsrv"
)

func testServer(t *testing.T) *systemServer {
	srv, _ := tsrv.RPCServer(t, false)
	srv.Config().Signup.AutoConfirm = true
	return newSystemServer(srv)
}
