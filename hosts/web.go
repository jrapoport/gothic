package hosts

import (
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/hosts/rpc/account"
	"github.com/jrapoport/gothic/hosts/rpc/health"
	"github.com/jrapoport/gothic/hosts/rpc/user"
)

const rpcWebName = "rpc-web"

// NewRPCWebHost creates a new rpc host.
func NewRPCWebHost(a *core.API, address string) core.Hosted {
	return rpc.NewHost(a, rpcWebName, address,
		[]rpc.RegisterServer{
			account.RegisterServer,
			health.RegisterServer,
			user.RegisterServer,
		}, rpc.Authentication())
}
