package hosts

import (
	"net"

	"github.com/jrapoport/gothic/api"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// RegisterRpcServer is the function prototype for registering an RPC server.
type RegisterRpcServer func(s *grpc.Server, srv *RPCHost)

// RPCHost represents a gRPC host.
type RPCHost struct {
	*api.API
	*logrus.Entry
	hostAndPort string
	servers     []RegisterRpcServer
}

// NewRpcHost creates a new RPCHost.
func NewRpcHost(a *api.API, name string, hostAndPort string, servers []RegisterRpcServer) *RPCHost {
	log := logrus.WithField("server", name)
	return &RPCHost{a, log, hostAndPort, servers}
}

func (h *RPCHost) ListenAndServe(opts ...grpc.ServerOption) {
	lis, err := net.Listen("tcp", h.hostAndPort)
	if err != nil {
		h.WithError(err).Fatal("rpc server listen failed")
	}

	server := grpc.NewServer(opts...)
	for _, s := range h.servers {
		s(server, h)
	}

	if conf.Debug {
		// Register reflection service on gRPC server.
		reflection.Register(server)
	}

	done := make(chan struct{})
	defer close(done)
	go func() {
		utils.WaitForTermination(h, done)
		h.Info("shutting down rpc server...")
		server.GracefulStop()
	}()

	if err := server.Serve(lis); err != nil {
		h.WithError(err).Fatal("rpc server failed to start")
	}
}

// RpcError wraps an rpc error code.
func (h *RPCHost) RpcError(c codes.Code, msg string) error {
	err := status.Error(c, msg)
	h.Error(err)
	return err
}

// RpcErrorf wraps a formatted rpc error code.
func (h *RPCHost) RpcErrorf(c codes.Code, format string, a ...interface{}) error {
	err := status.Errorf(c, format, a...)
	h.Error(err)
	return err
}
