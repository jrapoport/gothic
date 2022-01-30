package rpc

import (
	"context"
	"net"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/jrapoport/gothic/core"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RegisterServer is the function prototype for registering an rpc server.
type RegisterServer func(s *grpc.Server, srv *Server)

// Host represents a gRPC host.
type Host struct {
	*core.Host
	server *grpc.Server
}

type authOption struct {
	grpc.EmptyServerOption
}

// Authentication option to enable jwt authentication.
func Authentication() grpc.ServerOption {
	return authOption{}
}

// NewHost creates a new Host.
func NewHost(a *core.API, name string, address string, reg []RegisterServer, opt ...grpc.ServerOption) *Host {
	h := core.NewHost(a, name, address)
	h.Logger = h.WithName("grpc")
	l := grpcLogger(h.Config().Level)
	unary := []grpc.UnaryServerInterceptor{
		grpc_logrus.UnaryServerInterceptor(l, grpc_logrus.WithDecider(func(fullMethodName string, err error) bool {
			switch fullMethodName {
			case "/grpc.health.v1.Health/Check":
				return false
			case "/grpc.health.v1.Health/Watch":
				return false
			default:
				return true
			}
		})),
		grpc_recovery.UnaryServerInterceptor(),
		// ratelimit.UnaryServerInterceptor(),
	}
	stream := []grpc.StreamServerInterceptor{
		grpc_logrus.StreamServerInterceptor(l),
		grpc_recovery.StreamServerInterceptor(),
		// ratelimit.StreamServerInterceptor(),
	}
	if h.Config().Tracer.Enabled {
		unary = append(unary, grpc_opentracing.UnaryServerInterceptor())
		stream = append(stream, grpc_opentracing.StreamServerInterceptor())
	}

	for _, o := range opt {
		switch o.(type) {
		case authOption:
			// ignore this so we don't break reflection
			//if s.Config().IsDebug() {
			//	break
			//}
			auth := NewAuthenticator(h.Config().JWT, h.Log())
			unary = append(unary, auth.UnaryServerInterceptor())
			stream = append(stream, auth.StreamServerInterceptor())
		}
	}
	if len(unary) > 0 {
		ci := grpc.ChainUnaryInterceptor(unary...)
		opt = append(opt, ci)
	}
	if len(stream) > 0 {
		ci := grpc.ChainStreamInterceptor(stream...)
		opt = append(opt, ci)
	}
	server := grpc.NewServer(opt...)
	for _, r := range reg {
		srv := NewServer(h.Server.Clone())
		r(server, srv)
	}
	return &Host{h, server}
}

// ListenAndServe starts the rpc server.
func (h *Host) ListenAndServe() error {
	dbg := h.Config().IsDebug()
	h.Start(func(lis net.Listener) error {
		if dbg {
			// Register reflection service on gRPC server.
			reflection.Register(h.server)
		}
		return h.server.Serve(lis)
	})
	return h.Host.ListenAndServe()
}

// Shutdown stops the rpc server.
func (h *Host) Shutdown() error {
	h.Stop(func(context.Context) error {
		h.server.GracefulStop()
		return nil
	})
	return h.Host.Shutdown()
}

func grpcLogger(level string) *logrus.Entry {
	lvl, _ := logrus.ParseLevel(level)
	l := logrus.New()
	l.SetLevel(lvl)
	return l.WithField("protocol", "grpc")
}
