package user

// Role is te user role.
type Role int8

const (
	// InvalidRole is an invalid role.
	InvalidRole Role = iota - 1
	// RoleSystem is a system role.
	RoleSystem
	// RoleUser is a user role.
	RoleUser
	// RoleAdmin is an admin role.
	RoleAdmin
	// RoleSuper is a super admin role.
	RoleSuper
)

// Valid returns true if the provider is valid.
func (r Role) Valid() bool {
	switch r {
	case RoleUser:
		break
	case RoleAdmin:
		break
	case RoleSuper:
		break
	case RoleSystem:
		break
	default:
		return false
	}
	return true
}

// ToRole returns a role for a string.
func ToRole(s string) Role {
	switch s {
	case "system":
		return RoleSystem
	case "user":
		return RoleUser
	case "admin":
		return RoleAdmin
	case "super":
		return RoleSuper
	default:
		return InvalidRole
	}
}

func (r Role) String() string {
	switch r {
	case RoleSystem:
		return "system"
	case RoleUser:
		return "user"
	case RoleAdmin:
		return "admin"
	case RoleSuper:
		return "super"
	default:
		return ""
	}
}

// MarshalJSON marshals a role to json
func (r Role) MarshalJSON() ([]byte, error) {
	return []byte("\"" + r.String() + "\""), nil
}
