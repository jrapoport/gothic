package user

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/utils"
	"gorm.io/gorm"
)

func init() {
	var userIndexes = []string{
		"idx_email",
	}
	store.AddAutoMigrationWithIndexes("5000-users",
		User{}, userIndexes)
}

// Status is the user status
type Status int8

const (
	// Invalid user status
	Invalid Status = iota - 1
	// Banned user status
	Banned
	// Locked user status
	Locked
	// Restricted user status
	Restricted
	// Active user status
	Active
	// Verified user status
	Verified
)

// User represents a registered user with email/password authentication
// TODO: support additional verification (beyond email confirmation) via VerifiedAt
type User struct {
	ID          uuid.UUID        `json:"id" gorm:"primaryKey;type:char(36)"`
	Provider    provider.Name    `json:"provider" gorm:"type:varchar(255)"`
	Role        Role             `json:"role" gorm:"type:varchar(36)"`
	Status      Status           `json:"status"`
	Email       string           `json:"email" gorm:"uniqueIndex;type:varchar(320)"`
	Username    string           `json:"username" gorm:"type:varchar(255)"`
	Password    []byte           `json:"-" gorm:"type:binary(60)"`
	Data        types.Map        `json:"data"`
	Metadata    types.Map        `json:"metadata"`
	SignupCode  *uint            `json:"signup_code"`
	Linked      []account.Linked `json:"linked"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	LoginAt     *time.Time       `json:"login_at,omitempty"`
	ConfirmedAt *time.Time       `json:"confirmed_at,omitempty"`
	VerifiedAt  *time.Time       `json:"verified_at,omitempty"`
	InvitedAt   *time.Time       `json:"invited_at,omitempty"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at"`
}

// NewUser initializes a new user from an email, password and user data.
func NewUser(p provider.Name, role Role, email, username string, password []byte, data, meta types.Map) *User {
	if !role.Valid() || role == RoleSystem {
		return nil
	} else if email == "" {
		return nil
	} else if password == nil {
		return nil
	}
	u := &User{
		ID:       uuid.New(),
		Provider: p,
		Role:     role,
		Status:   Restricted,
		Email:    email,
		Username: username,
		Password: password,
		Data:     data,
		Metadata: meta,
	}
	return u
}

// BeforeSave runs before create or update.
func (u *User) BeforeSave(*gorm.DB) error {
	if u.IsSystemUser() {
		return errors.New("invalid user id (system)")
	}
	if u.Provider == provider.Unknown {
		return errors.New("invalid provider")
	}
	return nil
}

// EmailAddress returns the email address for the user.
func (u User) EmailAddress() *mail.Address {
	return &mail.Address{
		Name:    u.Username,
		Address: u.Email,
	}
}

// Authenticate returns nil if the password matches.
func (u User) Authenticate(pw string) error {
	if u.IsLocked() {
		return errors.New("invalid user")
	}
	return utils.CheckPassword(u.Password, pw)
}

// IsSystemUser returns true if a user is a system account.
func (u *User) IsSystemUser() bool {
	return u.ID == SystemID
}

// Valid returns true if the user is valid.
func (u User) Valid() bool {
	if u.CreatedAt.IsZero() {
		return false
	}
	if u.DeletedAt.Valid {
		return false
	}
	return !u.IsSystemUser()
}

// IsAdmin returns true if a user is an admin.
func (u User) IsAdmin() bool {
	return u.Valid() && u.Role >= RoleAdmin
}

// IsBanned returns true if the user is not banned.
func (u User) IsBanned() bool {
	return !u.Valid() || u.Status <= Banned
}

// IsLocked returns true if the user is not locked or banned.
func (u User) IsLocked() bool {
	return u.IsBanned() || u.Status <= Locked
}

// IsRestricted returns true if the user is not confirmed.
func (u User) IsRestricted() bool {
	return !u.IsLocked() && (u.Status <= Restricted || u.ConfirmedAt == nil)
}

// IsConfirmed returns true if a user has confirmed their email address.
func (u User) IsConfirmed() bool {
	return !u.IsLocked() && u.ConfirmedAt != nil
}

// IsActive returns true if a user has confirmed their email address.
func (u User) IsActive() bool {
	return u.IsConfirmed() && u.Status >= Active
}

// IsVerified returns true  if a user has been verified.
func (u User) IsVerified() bool {
	return u.IsConfirmed() && u.VerifiedAt != nil && u.Status >= Verified
}
