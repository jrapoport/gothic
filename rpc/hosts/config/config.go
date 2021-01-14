package config

//go:generate protoc -I=. --go_out=plugins=grpc:. --go_opt=paths=source_relative config.proto

import (
	"context"

	"github.com/jrapoport/gothic/rpc/hosts"
	"github.com/segmentio/encoding/json"
	"google.golang.org/grpc/codes"
)

type rpcConfigHost struct {
	*hosts.RpcHost
}

var _ ConfigServer = (*rpcConfigHost)(nil)

func NewConfigHost(h *hosts.RpcHost) *rpcConfigHost {
	return &rpcConfigHost{h}
}

func (r *rpcConfigHost) Settings(_ context.Context,
	_ *SettingsRequest) (*SettingsResponse, error) {
	s := r.API.Settings()
	b, err := json.Marshal(s)
	if err != nil {
		return nil, r.RpcErrorf(codes.Internal, "%v", err)
	}
	res := SettingsResponse{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, r.RpcErrorf(codes.Internal, "%v", err)
	}
	return &res, nil
}
