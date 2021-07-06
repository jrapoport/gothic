package events

// Event dispatch event
type Event string

const (
	// Unknown event
	Unknown Event = ""
	// Confirmed event
	Confirmed Event = "confirmed"
	// Login event
	Login Event = "login"
	// Logout event
	Logout Event = "logout"
	// Signup event
	Signup Event = "signup"
	// All must be last
	All Event = "all"
)
