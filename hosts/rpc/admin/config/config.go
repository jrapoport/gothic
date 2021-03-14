package config

//go:generate protoc -I=. --go_out=plugins=grpc:. --go_opt=paths=source_relative config.proto

import (
	"context"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/segmentio/encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type configServer struct {
	*rpc.Server
}

var _ ConfigServer = (*configServer)(nil)

func newConfigServer(srv *rpc.Server) *configServer {
	srv.FieldLogger = srv.WithField("module", "config")
	return &configServer{srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	RegisterConfigServer(s, newConfigServer(srv))
}

// Settings returns the settings for a server.
func (s *configServer) Settings(_ context.Context, _ *SettingsRequest) (*SettingsResponse, error) {
	set := s.API.Settings()
	b, err := json.Marshal(set)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	res := SettingsResponse{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	return &res, nil
}
