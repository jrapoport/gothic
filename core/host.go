package core

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Hosted is a host interface
type Hosted interface {
	Name() string
	Address() string
	Online() bool
	ListenAndServe() error
	Shutdown() error
}

// ServeFunc is a function prototype for start.
type ServeFunc func(net.Listener) error

// ShutdownFunc is a function prototype for shutdown.
type ShutdownFunc func(context.Context) error

// Host generalized core server host.
type Host struct {
	Server
	lis    net.Listener
	name   string
	addr   string
	wait   int32
	online bool
	start  ServeFunc
	stop   ShutdownFunc
	mu     sync.RWMutex
}

var _ Hosted = (*Host)(nil)

// NewHost creates a new Host.
func NewHost(a *API, name, address string) *Host {
	s := *NewServer(a, name)
	if address == "" {
		s.Warnf("%s host is disabled", name)
	}
	return &Host{
		Server: s,
		name:   name,
		addr:   address,
	}
}

// Name returns the name of the host.
func (s *Host) Name() string {
	return s.name
}

// Address returns the listening address of the host.
func (s *Host) Address() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.lis == nil {
		return s.addr
	}
	return s.lis.Addr().String()
}

// Online returns true if the host is online.
func (s *Host) Online() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.online
}

// Start starts a core server host.
func (s *Host) Start(start ServeFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.start = start
}

// ListenAndServe listens on a port and serves.
func (s *Host) ListenAndServe() error {
	if s.start == nil || s.addr == "" {
		return nil
	}
	started := sync.WaitGroup{}
	defer started.Wait()
	s.Infof("%s server %s starting...", s.name, s.addr)
	if err := s.startListing(); err != nil {
		return err
	}
	started.Add(1)
	go func() {
		defer func() {
			err := s.stopListening()
			if err != nil {
				s.Error(err)
			}
		}()
		s.Infof("%s server %s started", s.name, s.addr)
		s.mu.Lock()
		s.online = true
		s.mu.Unlock()
		started.Done()
		if err := s.start(s.lis); err != nil {
			err = fmt.Errorf("%s server %s failed to start: %w", s.name, s.addr, err)
			s.Error(err)
			return
		}
		s.mu.Lock()
		s.online = false
		s.mu.Unlock()
		s.Infof("%s server %s stopped", s.name, s.addr)
	}()
	return nil
}

// Stop stops a core server host.
func (s *Host) Stop(stop ShutdownFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stop = stop
}

// Shutdown shuts down the host.
func (s *Host) Shutdown() error {
	if s.stop == nil || s.addr == "" {
		return s.stopListening()
	}
	timeout := time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.Infof("%s server %s shutting down...", s.name, s.addr)
	if err := s.stop(ctx); err != nil {
		s.Error(err)
		return err
	}
	s.Infof("%s server %s shut down", s.name, s.addr)
	waitWithCtx(ctx, &s.wait)
	return nil
}

func (s *Host) startListing() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// create a listener on TCP port
	s.lis, err = net.Listen("tcp", s.addr)
	if err != nil {
		err = fmt.Errorf("%s failed to listen %s: %w", s.name, s.addr, err)
		return err
	}
	s.Infof("%s server listening on %s", s.name, s.addr)
	// run this instead of wg.Add(1)
	atomic.AddInt32(&s.wait, 1)
	return
}

func (s *Host) stopListening() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.lis == nil {
		return nil
	}
	if err := s.lis.Close(); err != nil {
		msg := err.Error()
		if !strings.Contains(msg, "use of closed network connection") {
			s.Error(err)
			return err
		}
	}
	s.Infof("%s server stopped listening on %s ", s.name, s.addr)
	// run this instead of wg.Done()
	atomic.AddInt32(&s.wait, -1)
	return nil
}

// waitWithCtx returns when passed counter drops to zero or when context is cancelled
func waitWithCtx(ctx context.Context, counter *int32) {
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if atomic.LoadInt32(counter) == 0 {
				return
			}
		}
	}
}
