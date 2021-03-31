package user

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/utils"
)

// SuperAdminID is the user id for the super admin account.
var SuperAdminID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

// NewSuperAdmin returns a new super admin user.
func NewSuperAdmin(password string) *User {
	return &User{
		ID:       SuperAdminID,
		Role:     RoleSuper,
		Status:   Verified,
		Password: utils.HashPassword(password),
	}
}
