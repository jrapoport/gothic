package auditlog

// Type is the type of the log entry
type Type int

const (
	//Unknown log entry type
	Unknown Type = iota - 1
	//System log entry type
	System
	//Account log entry type
	Account
	//Token log entry type
	Token
	//User log entry type
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
