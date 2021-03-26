package settings

import (
	"context"
	"github.com/jrapoport/gothic/api/grpc/rpc/admin/settings"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/segmentio/encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type settingsServer struct {
	settings.UnimplementedSettingsServer
	*rpc.Server
}

var _ settings.SettingsServer = (*settingsServer)(nil)

func newSettingsServer(srv *rpc.Server) *settingsServer {
	srv.FieldLogger = srv.WithField("module", "config")
	return &settingsServer{Server: srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	settings.RegisterSettingsServer(s, newSettingsServer(srv))
}

// Settings returns the settings for a server.
func (s *settingsServer) Settings(_ context.Context, _ *settings.SettingsRequest) (*settings.SettingsResponse, error) {
	set := s.API.Settings()
	b, err := json.Marshal(set)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	res := settings.SettingsResponse{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	return &res, nil
}
