package hosts

import (
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/admin"
	"github.com/jrapoport/gothic/hosts/rpc/health"
)

const rpcName = "rpc"

// NewRPCHost creates a new rpc host.
func NewRPCHost(a *core.API, address string) core.Hosted {
	return rpc.NewHost(a, rpcName, address,
		[]rpc.RegisterServer{
			admin.RegisterServer,
			health.RegisterServer,
		})
}
