package config

// Invites config
type Invites string

// Invite permissions
const (
	Disabled Invites = ""
	Users    Invites = "user"
	Admins   Invites = "admin"
	Super    Invites = "super"
)

// Signup config
type Signup struct {
	// Disabled disables all user signups.
	Disabled bool `json:"disabled"`
	// AutoConfirm automatically confirms a user.
	AutoConfirm bool `json:"autoconfirm"`
	// Code if true, require a signup code.
	Code bool `json:"code"`
	// Invites controls access to invitations.
	Invites Invites `json:"invites"`
	// Username if true, require a username.
	Username bool `json:"username"`
	// Default enables which default values will be supplied if absent.
	Default SignupDefaults `json:"default"`
}

// CanSendInvites returns true if invites are enabled.
func (s Signup) CanSendInvites() bool {
	if s.Disabled {
		return false
	}
	switch s.Invites {
	case Users, Admins, Super:
		return true
	default:
		return false
	}
}

// SignupDefaults config
type SignupDefaults struct {
	// Username generates a random username if one is not present.
	Username bool `json:"username"`
	// Color generates a random user color if one is not present.
	Color bool `json:"color"`
}
