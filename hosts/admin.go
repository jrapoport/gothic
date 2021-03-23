package hosts

import (
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/admin"
)

const adminName = "admin"

// NewAdminHost creates a new admin host.
func NewAdminHost(a *core.API, address string) core.Hosted {
	return rpc.NewHost(a, adminName, address,
		[]rpc.RegisterServer{
			admin.RegisterServer,
		})
}
