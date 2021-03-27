package admin

import (
	"context"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/segmentio/encoding/json"
	"google.golang.org/grpc/codes"
)

// Settings returns the settings for a server.
func (s *adminServer) Settings(_ context.Context, _ *admin.SettingsRequest) (*admin.SettingsResponse, error) {
	set := s.API.Settings()
	b, err := json.Marshal(set)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	res := admin.SettingsResponse{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	return &res, nil
}
