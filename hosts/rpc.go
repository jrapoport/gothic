package hosts

import (
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/health"
	"github.com/jrapoport/gothic/hosts/rpc/system"
)

const rpcName = "rpc"

// NewRPCHost creates a new rpc host.
func NewRPCHost(a *core.API, address string) core.Hosted {
	return rpc.NewHost(a, rpcName, address,
		[]rpc.RegisterServer{
			system.RegisterServer,
			health.RegisterServer,
		})
}
