package auditlog

// Type is the type of the log entry
type Type int

// Audit log entry types
const (
	Unknown Type = iota - 1
	System
	Account
	Token
	User
)

// TypeFromString returns the log type for a string.
func TypeFromString(s string) Type {
	switch s {
	case "system":
		return System
	case "account":
		return Account
	case "token":
		return Token
	case "user":
		return User
	default:
		return Unknown
	}
}

func (t Type) String() string {
	switch t {
	case System:
		return "system"
	case Account:
		return "account"
	case Token:
		return "token"
	case User:
		return "user"
	default:
		return ""
	}
}
