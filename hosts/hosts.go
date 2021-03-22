package hosts

import (
	"fmt"
	"sync"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts/health"
)

type hostMan struct {
	hosted  map[string]core.Hosted
	mu      sync.RWMutex
	running bool
}

var hm *hostMan

func init() {
	hm = &hostMan{hosted: map[string]core.Hosted{}}
}

// Start starts the configured hosts.
func Start(a *core.API, c *config.Config) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hosts := []struct {
		name   string
		addr   *string
		hostFn func(a *core.API, addr string) core.Hosted
	}{
		{adminName, &c.Admin, NewAdminHost},
		{rpcName, &c.RPC, NewRPCHost},
		{rpcWebName, &c.RPCWeb, NewRPCWebHost},
		{restName, &c.REST, NewRESTHost},
		{health.Name, &c.Health, func(a *core.API, addr string) core.Hosted {
			return health.NewHealthHost(a, addr, hm.hosted)
		}},
	}
	for _, h := range hosts {
		s := h.hostFn(a, *h.addr)
		if s == nil {
			err := fmt.Errorf("unable to load %s host", h.name)
			return err
		}
		err := s.ListenAndServe()
		if err != nil {
			return err
		}
		hm.hosted[h.name] = s
		*h.addr = s.Address()
	}
	hm.running = true
	l := c.Log().WithField("logger", "grpc")
	grpc_logrus.ReplaceGrpcLogger(l)
	return nil
}

// Running returns true if the manager is running.
func Running() bool {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	return hm.running
}

// Shutdown stops the configured hosts.
func Shutdown() error {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	var err error
	for _, h := range hm.hosted {
		e := h.Shutdown()
		if e != nil {
			err = fmt.Errorf("%w", err)
		}
	}
	hm.running = false
	return err
}
