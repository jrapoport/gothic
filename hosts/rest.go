package hosts

import (
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts/rest"
	"github.com/jrapoport/gothic/hosts/rest/account"
	"github.com/jrapoport/gothic/hosts/rest/admin"
	"github.com/jrapoport/gothic/hosts/rest/health"
	"github.com/jrapoport/gothic/hosts/rest/user"
)

const restName = "rest"

// NewRESTHost creates a new rest host.
func NewRESTHost(a *core.API, address string) core.Hosted {
	return rest.NewHost(a, restName, address,
		[]rest.RegisterServer{
			admin.RegisterServer,
			account.RegisterServer,
			health.RegisterServer,
			user.RegisterServer,
		})
}
