package core

import "github.com/jrapoport/gothic/core/health"

// HealthCheck returns the health check.
func (a *API) HealthCheck() health.Health {
	return health.Check(a.config)
}
