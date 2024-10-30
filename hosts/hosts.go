package hosts

import (
	"fmt"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"sync"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts/health"
	"github.com/sirupsen/logrus"
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
	grpcLogLevel(c.Level)
	hosts := []struct {
		name   string
		addr   *string
		hostFn func(a *core.API, addr string) core.Hosted
	}{
		{adminName, &c.AdminAddress, NewAdminHost},
		{rpcName, &c.RPCAddress, NewRPCHost},
		{rpcWebName, &c.RPCWebAddress, NewRPCWebHost},
		{restName, &c.RESTAddress, NewRESTHost},
		{health.Name, &c.HealthAddress, func(a *core.API, addr string) core.Hosted {
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

func grpcLogLevel(level string) {
	lvl, _ := logrus.ParseLevel(level)
	l := logrus.New()
	l.SetLevel(lvl)
	grpc_logrus.ReplaceGrpcLogger(l.WithField("protocol", "grpc"))
}
