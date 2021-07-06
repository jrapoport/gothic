package admin

//go:generate protoc -I=. -I=.. --go_out=plugins=grpc:. --go_opt=paths=source_relative admin.proto

import (
	"context"

	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/segmentio/encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type adminServer struct {
	*rpc.Server
}

var _ AdminServer = (*adminServer)(nil)

func newAdminServer(srv *rpc.Server) *adminServer {
	srv.FieldLogger = srv.WithField("module", "admin")
	return &adminServer{srv}
}

// RegisterServer registers a new admin server.
func RegisterServer(s *grpc.Server, srv *rpc.Server) {
	RegisterAdminServer(s, newAdminServer(srv))
}

// Settings returns the settings for a server.
func (s *adminServer) Settings(_ context.Context, _ *SettingsRequest) (*SettingsResponse, error) {
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
