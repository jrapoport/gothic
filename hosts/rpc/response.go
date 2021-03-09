package rpc

//go:generate protoc -I=. --go_out=plugins=grpc:. --go_opt=paths=source_relative response.proto

import (
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/utils"
	"google.golang.org/protobuf/types/known/structpb"
)

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

// MaskEmail masks the email field of the user response
func (x *UserResponse) MaskEmail() {
	x.Email = utils.MaskEmail(x.Email)
}

// NewBearerResponse returns a BearerResponse from a BearerToken
func NewBearerResponse(bt *tokens.BearerToken) *BearerResponse {
	return &BearerResponse{
		Type:    bt.Class().String(),
		Access:  bt.String(),
		Refresh: bt.RefreshToken.String(),
	}
}
