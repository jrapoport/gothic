package settings

//go:generate protoc -I=. --go_out=plugins=grpc:. --go_opt=paths=source_relative settings.proto

import (
	"context"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/segmentio/encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type settingsServer struct {
	*rpc.Server
}

var _ SettingsServer = (*settingsServer)(nil)

func newSettingsServer(srv *rpc.Server) *settingsServer {
	srv.FieldLogger = srv.WithField("module", "config")
	return &settingsServer{srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	RegisterSettingsServer(s, newSettingsServer(srv))
}

// Settings returns the settings for a server.
func (s *settingsServer) Settings(_ context.Context, _ *SettingsRequest) (*SettingsResponse, error) {
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
