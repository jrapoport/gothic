package rpc

import (
	"github.com/jrapoport/gothic/api/grpc/rpc"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/utils"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserResponse api.UserResponse

// NewUserResponse returns a UserResponse for the supplied user.
func NewUserResponse(u *user.User) (*UserResponse, error) {
	data, err := structpb.NewStruct(u.Data)
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		Role:     u.Role.String(),
		Email:    u.Email,
		Username: u.Username,
		Data:     data,
	}, nil
}

// MaskEmail masks Personally Identifiable Information (PII) from the user response
func (r *UserResponse) MaskEmail() {
	r.Email = utils.MaskEmail(r.Email)
}

// NewBearerResponse returns a BearerResponse from a BearerToken
func NewBearerResponse(bt *tokens.BearerToken) *api.BearerResponse {
	res := &api.BearerResponse{
		Type:    bt.Class().String(),
		Access:  bt.String(),
		Refresh: bt.RefreshToken.String(),
	}
	if bt.ExpiredAt != nil {
		res.ExpiresAt = timestamppb.New(*bt.ExpiredAt)
	}
	return res
}
