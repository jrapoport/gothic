package auditlog

// Account actions
const (
	// Banned action
	Banned Action = "banned"
	// CodeSent action
	CodeSent Action = "code_sent"
	// ConfirmSent action
	ConfirmSent Action = "confirm_sent"
	// Confirmed action
	Confirmed Action = "confirmed"
	// Deleted action
	Deleted Action = "deleted"
	// Linked action
	Linked Action = "linked"
	// Signup action
	Signup Action = "signup"
)

// System actions
const (
	// Startup action
	Startup Action = "startup"
	// Shutdown action
	Shutdown Action = "shutdown"
)

// Token actions
const (
	// Granted action
	Granted Action = "granted"
	// Refreshed action
	Refreshed Action = "refreshed"
	// Revoked action
	Revoked Action = "revoked"
	// RevokedAll action
	RevokedAll Action = "revoked_all"
)

// User actions
const (
	// ChangeRole action
	ChangeRole Action = "change_role"
	// Email action
	Email Action = "email"
	// Login action
	Login Action = "login"
	// Logout action
	Logout Action = "logout"
	// Password action
	Password Action = "password"
	// Updated action
	Updated Action = "updated"
)

// Action is action captured by the log entry.
type Action string

// Type returns the type for the action.
func (a Action) Type() Type {
	switch a {
	// System actions
	case Startup:
		return System
	case Shutdown:
		return System
	// Account actions
	case Signup:
		return Account
	case CodeSent:
		return Account
	case ConfirmSent:
		return Account
	case Confirmed:
		return Account
	case Linked:
		return Account
	case Banned:
		return Account
	case Deleted:
		return Account
	// Token actions
	case Granted:
		return Token
	case Refreshed:
		return Token
	case Revoked:
		return Token
	case RevokedAll:
		return Token
	// User actions
	case Login:
		return User
	case Logout:
		return User
	case Password:
		return User
	case Email:
		return User
	case Updated:
		return User
	case ChangeRole:
		return User
	// Unknown action
	default:
		return Unknown
	}
}

func (a Action) String() string {
	return string(a)
}
