package admin

import (
	"context"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/segmentio/encoding/json"
	"google.golang.org/grpc/codes"
)

// Settings returns the settings for a server.
func (s *server) Settings(ctx context.Context, _ *admin.SettingsRequest) (*admin.SettingsResponse, error) {
	_, err := s.adminRequestContext(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	set := s.API.Settings()
	// we can safely ignore this error
	// bc we tightly control the type
	b, _ := json.Marshal(set)
	res := admin.SettingsResponse{}
	// we can safely ignore this error
	// bc we tightly control the type
	_ = json.Unmarshal(b, &res)
	return &res, nil
}
