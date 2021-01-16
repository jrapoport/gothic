package web

import (
	"github.com/jrapoport/gothic/api"
	"github.com/jrapoport/gothic/rpc/hosts"
	"github.com/jrapoport/gothic/rpc/hosts/health"
	"google.golang.org/grpc"
)

type rpcWebServer struct {
	*hosts.RPCHost
}

// NewRPCWebServer creates a new gRPC-Web server.
func NewRPCWebServer(a *api.API, hostAndPort string) *rpcWebServer {
	s := hosts.NewRpcHost(a, "rpc-web", hostAndPort, []hosts.RegisterRpcServer{
		func(s *grpc.Server, srv *hosts.RPCHost) {
			hs := health.NewHealthHost(srv)
			health.RegisterHealthServer(s, hs)
		},
	})
	return &rpcWebServer{s}
}
