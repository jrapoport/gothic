package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Token is the interface for tokens.
type Token interface {
	Class() Class
	Usage() Type
	IssuedTo() uuid.UUID
	Issued() time.Time
	LastUsed() time.Time
	ExpirationDate() time.Time
	Revoked() time.Time
	Usable() bool
	Use() // make this UseToken

	fmt.Stringer
}

// Type is the type of code (e.g. single use)
type Type uint8

const (
	// Infinite indicates this code may be used an infinite number of times.
	Infinite Type = iota
	// Single is a token that can be used once.
	Single
	// Multi is a token that can be used repeatedly.
	Multi
	// Timed is a token that must be used before it expires.
	Timed
)

func (t Type) String() string {
	switch t {
	case Infinite:
		return "infinite"
	case Single:
		return "single"
	case Multi:
		return "multi"
	case Timed:
		return "timed"
	default:
		return "invalid"
	}
}

// Class is the class of the token.
type Class string

func (c Class) String() string {
	return string(c)
}

const (
	// Access is a generic access token
	Access Class = "access"
	// Confirm is an email confirmation token.
	Confirm Class = "confirm"
	// Refresh is a refresh token.
	Refresh Class = "refresh"
	// Auth is an authorization token.
	Auth Class = "auth"
)
