package user

import "github.com/google/uuid"

// SystemID is the user id for system accounts.
var SystemID = uuid.Nil

// NewSystemUser returns a new system user.
func NewSystemUser() *User {
	return &User{
		ID:   SystemID,
		Role: RoleSystem,
	}
}
