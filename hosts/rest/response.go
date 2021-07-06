package rest

import (
	"time"

	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/utils"
)

// UserResponse contains an http user response.
type UserResponse struct {
	UserID   string          `json:"user_id,omitempty"`
	Role     string          `json:"role"`
	Email    string          `json:"email,omitempty"`
	Username string          `json:"username,omitempty"`
	Data     types.Map       `json:"data"`
	Token    *BearerResponse `json:"token,omitempty"`
}

// NewUserResponse returns a UserResponse for the supplied user.
func NewUserResponse(u *user.User) *UserResponse {
	return &UserResponse{
		UserID:   u.ID.String(),
		Role:     u.Role.String(),
		Email:    u.Email,
		Username: u.Username,
		Data:     u.Data,
	}
}

// MaskEmail masks the email field of the user response
func (r *UserResponse) MaskEmail() {
	r.Email = utils.MaskEmail(r.Email)
}

// BearerResponse is the
type BearerResponse struct {
	Type      string     `json:"type"`
	Access    string     `json:"access,omitempty"`
	Refresh   string     `json:"refresh,omitempty"`
	ID        string     `json:"id,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// NewBearerResponse returns a BearerResponse from a BearerToken
func NewBearerResponse(bt *tokens.BearerToken) *BearerResponse {
	return &BearerResponse{
		Type:      bt.Class().String(),
		Access:    bt.String(),
		Refresh:   bt.RefreshToken.String(),
		ExpiresAt: bt.ExpiredAt,
	}
}
