package health

//go:generate protoc -I=. --go_out=plugins=grpc:. --go_opt=paths=source_relative health.proto

import (
	"context"

	"github.com/jrapoport/gothic/rpc/hosts"
)

type rpcHealthHost struct {
	*hosts.RPCHost
}

var _ HealthServer = (*rpcHealthHost)(nil)

func NewHealthHost(h *hosts.RPCHost) *rpcHealthHost {
	return &rpcHealthHost{h}
}

func (r *rpcHealthHost) HealthCheck(_ context.Context,
	_ *HealthCheckRequest) (*HealthCheckResponse, error) {
	hc := r.API.HealthCheck()
	return &HealthCheckResponse{
		Version:     hc["version"],
		Name:        hc["name"],
		Description: hc["description"],
	}, nil
}
