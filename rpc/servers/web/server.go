package web

import (
	"github.com/jrapoport/gothic/api"
	"github.com/jrapoport/gothic/rpc/hosts"
	"github.com/jrapoport/gothic/rpc/hosts/health"
	"google.golang.org/grpc"
)

type rpcWebServer struct {
	*hosts.RpcHost
}

// NewRpcWebServer creates a new gRPC-Web server.
func NewRpcWebServer(a *api.API, hostAndPort string) *rpcWebServer {
	s := hosts.NewRpcHost(a, "rpc-web", hostAndPort, []hosts.RegisterRpcServer{
		func(s *grpc.Server, srv *hosts.RpcHost) {
			hs := health.NewHealthHost(srv)
			health.RegisterHealthServer(s, hs)
		},
	})
	return &rpcWebServer{s}
}
