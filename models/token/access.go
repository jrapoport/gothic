package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/user"
	"gorm.io/gorm"
)

const (
	// InfiniteUse indicates this code may be used an infinite
	// number of times. This implies type Infinite.
	InfiniteUse = -1

	// SingleUse indicates this code may be used once.
	// This implies type Single.
	SingleUse = 1
)

// NoExpiration indicates the token will not expire.
const NoExpiration time.Duration = 0

// AccessToken holds an access token.
type AccessToken struct {
	gorm.Model
	UserID     uuid.UUID     `json:"user_id" gorm:"index:idx_user_id;uniqueIndex:idx_user_id_token;type:char(36)"`
	Type       Usage         `json:"type"`
	Token      string        `json:"token" gorm:"index:idx_token;uniqueIndex:idx_user_id_token"`
	MaxUses    int           `json:"max_uses"`
	Used       int           `json:"used"`
	UsedAt     *time.Time    `json:"used_at,omitempty"`
	Expiration time.Duration `json:"expiration,omitempty"`
	ExpiredAt  *time.Time    `json:"expired_at,omitempty"`
	Data       types.Map     `json:"data,omitempty"`
}

// AccessTokenIndexes are the db indexes for the access token in the db.
var AccessTokenIndexes = []string{
	"idx_user_id",
	"idx_token",
	"idx_user_id_token",
}

// NewAccessToken generates a new token. The type is inferred from uses.
func NewAccessToken(token string, uses int, exp time.Duration) *AccessToken {
	if token == "" {
		return nil
	}
	t := Multi
	switch {
	case uses <= InfiniteUse:
		uses = InfiniteUse
		t = Infinite
		break
	case uses <= 1:
		t = Single
		uses = 1
	default:
		break
	}
	maxUses := uses
	switch {
	case exp < NoExpiration:
		exp = NoExpiration
	case exp > NoExpiration:
		t = Timed
	}
	return &AccessToken{
		UserID:     user.SystemID,
		Type:       t,
		Token:      token,
		MaxUses:    maxUses,
		Expiration: exp,
	}
}

var _ Token = (*AccessToken)(nil)

// BeforeCreate runs after create
func (t AccessToken) BeforeCreate(_ *gorm.DB) (err error) {
	if t.Token == "" {
		return errors.New("invalid token")
	}
	return nil
}

// AfterCreate runs after create
func (t AccessToken) AfterCreate(tx *gorm.DB) (err error) {
	if t.Expiration <= NoExpiration {
		return nil
	}
	exp := t.CreatedAt.Add(t.Expiration)
	m := tx.Statement.Model
	return tx.Model(m).Update("expired_at", exp).Error
}

// Class returns the class of the access token.
func (t AccessToken) Class() Class {
	return Access
}

// Usage returns usage for the token.
func (t AccessToken) Usage() Usage {
	return t.Type
}

// IssuedTo returns the owner of the token.
func (t AccessToken) IssuedTo() uuid.UUID {
	return t.UserID
}

// Issued returns the time the token was issued.
func (t AccessToken) Issued() time.Time {
	return t.CreatedAt
}

// LastUsed returns the last time the token was used.
func (t AccessToken) LastUsed() time.Time {
	if t.UsedAt == nil {
		return time.Time{}
	}
	return *t.UsedAt
}

// ExpirationDate returns the expiration date (if set).
func (t AccessToken) ExpirationDate() time.Time {
	if t.ExpiredAt == nil {
		return time.Time{}
	}
	return *t.ExpiredAt
}

// Revoked returns the time token was revoked (if set).
func (t AccessToken) Revoked() time.Time {
	return t.DeletedAt.Time
}

// Usable returns true if the token is usable.
func (t AccessToken) Usable() bool {
	if t.Class() == "" {
		return false
	}
	if t.Token == "" {
		return false
	}
	if t.DeletedAt.Valid {
		return false
	}
	if t.ExpiredAt != nil && t.ExpiredAt.Before(time.Now().UTC()) {
		return false
	}
	if t.Type == Infinite || t.MaxUses == InfiniteUse {
		return true
	}
	return t.Used < t.MaxUses
}

// Use use the token.
func (t *AccessToken) Use() {
	now := time.Now().UTC()
	t.Used++
	t.UsedAt = &now
}

func (t AccessToken) String() string {
	return t.Token
}
