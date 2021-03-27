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
	Type       Type          `json:"type"`
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
func (at AccessToken) BeforeCreate(_ *gorm.DB) (err error) {
	if at.Token == "" {
		return errors.New("invalid token")
	}
	return nil
}

// AfterCreate runs after create
func (at AccessToken) AfterCreate(tx *gorm.DB) (err error) {
	if at.Expiration <= NoExpiration {
		return nil
	}
	exp := at.CreatedAt.Add(at.Expiration)
	m := tx.Statement.Model
	return tx.Model(m).Update("expired_at", exp).Error
}

// Class returns the class of the access token.
func (at AccessToken) Class() Class {
	return Access
}

// Type returns usage for the token.
func (at AccessToken) Usage() Type {
	return at.Type
}

// IssuedTo returns the owner of the token.
func (at AccessToken) IssuedTo() uuid.UUID {
	return at.UserID
}

// Issued returns the time the token was issued.
func (at AccessToken) Issued() time.Time {
	return at.CreatedAt
}

// LastUsed returns the last time the token was used.
func (at AccessToken) LastUsed() time.Time {
	if at.UsedAt == nil {
		return time.Time{}
	}
	return *at.UsedAt
}

// ExpirationDate returns the expiration date (if set).
func (at AccessToken) ExpirationDate() time.Time {
	if at.ExpiredAt == nil {
		return time.Time{}
	}
	return *at.ExpiredAt
}

// Revoked returns the time token was revoked (if set).
func (at AccessToken) Revoked() time.Time {
	return at.DeletedAt.Time
}

// Usable returns true if the token is usable.
func (at AccessToken) Usable() bool {
	if at.Token == "" {
		return false
	}
	if at.DeletedAt.Valid {
		return false
	}
	if at.ExpiredAt != nil && at.ExpiredAt.Before(time.Now().UTC()) {
		return false
	}
	if at.Type == Infinite || at.MaxUses == InfiniteUse {
		return true
	}
	return at.Used < at.MaxUses
}

// Use use the token.
func (at *AccessToken) Use() {
	now := time.Now().UTC()
	at.Used++
	at.UsedAt = &now
}

func (at AccessToken) String() string {
	return at.Token
}
