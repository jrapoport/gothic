package rpc

import (
	"github.com/jrapoport/gothic/api"
	"github.com/jrapoport/gothic/rpc/hosts"
	"github.com/jrapoport/gothic/rpc/hosts/config"
	"github.com/jrapoport/gothic/rpc/hosts/health"
	"google.golang.org/grpc"
)

type rpcServer struct {
	*hosts.RPCHost
}

// NewRpcServer creates a new gRPC server.
func NewRpcServer(a *api.API, hostAndPort string) *rpcServer {
	s := hosts.NewRpcHost(a, "rpc", hostAndPort, []hosts.RegisterRpcServer{
		func(s *grpc.Server, srv *hosts.RPCHost) {
			ch := config.NewConfigHost(srv)
			config.RegisterConfigServer(s, ch)
		},
		func(s *grpc.Server, srv *hosts.RPCHost) {
			hh := health.NewHealthHost(srv)
			health.RegisterHealthServer(s, hh)
		},
	})
	return &rpcServer{s}
}
