package auditlog

// Account actions
const (
	Banned      Action = "banned"
	CodeSent    Action = "code_sent"
	ConfirmSent Action = "confirm_sent"
	Confirmed   Action = "confirmed"
	Deleted     Action = "deleted"
	Signup      Action = "signup"
)

// System actions
const (
	Startup  Action = "startup"
	Shutdown Action = "shutdown"
)

// Token actions
const (
	Granted    Action = "granted"
	Refreshed  Action = "refreshed"
	Revoked    Action = "revoked"
	RevokedAll Action = "revoked_all"
)

// User actions
const (
	ChangeRole Action = "change_role"
	Email      Action = "email"
	Linked     Action = "linked"
	Login      Action = "login"
	Logout     Action = "logout"
	Password   Action = "password"
	Updated    Action = "updated"
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
	case Linked:
		return User
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
