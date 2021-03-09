package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/store"
)

// Token is the interface for tokens.
type Token interface {
	Class() Class
	Usage() Usage
	IssuedTo() uuid.UUID
	Issued() time.Time
	LastUsed() time.Time
	ExpirationDate() time.Time
	Revoked() time.Time
	Usable() bool
	Use() // make this UseToken

	fmt.Stringer
}

// Usage is the type of code (e.g. single use)
type Usage uint8

const (
	// Infinite indicates this code may be used an infinite number of times.
	Infinite Usage = iota

	// Single is a token that can be used once.
	Single

	// Multi is a token that can be used repeatedly.
	Multi

	// Timed is a token that must be used before it expires.
	Timed
)

func (t Usage) String() string {
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

// UseToken burns a usable token
func UseToken(conn *store.Connection, t Token) error {
	if !t.Usable() {
		return errors.New("invalid")
	}
	return conn.Transaction(func(tx *store.Connection) error {
		t.Use()
		err := tx.Save(t).Error
		if err != nil {
			return err
		}
		if t.Usable() {
			return nil
		}
		return tx.Delete(t).Error
	})
}
