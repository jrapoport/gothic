package events

// Event dispatch event
type Event string

// Events
const (
	Unknown   Event = ""
	Confirmed Event = "confirmed"
	Login     Event = "login"
	Logout    Event = "logout"
	Signup    Event = "signup"
	All       Event = "all" // must be last
)
