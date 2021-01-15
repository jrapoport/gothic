package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/storage"
	"github.com/vcraescu/go-paginator/v2"
	"github.com/vcraescu/go-paginator/v2/adapter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	UserRole       Role = "user"
	AdminRole      Role = "admin"
	SuperAdminRole Role = "super-admin"
)

const SystemUserID = "0"

var SystemUserUUID = uuid.Nil

func init() {
	storage.AddMigration(&User{})
}

// User represents a registered user with email/password authentication
type User struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey"`

	Username          string `json:"username" gorm:"type:varchar(320)"`
	Email             string `json:"email" gorm:"type:varchar(320)"`
	EncryptedPassword string `json:"-"`

	// user must change their password
	ChangePassword bool `json:"change_password"`

	Role         string     `json:"role" gorm:"type:varchar(255)"`
	LastSignInAt *time.Time `json:"last_sign_in_at,omitempty"`

	ConfirmationToken  string     `json:"-" gorm:"type:varchar(255)"`
	ConfirmationSentAt *time.Time `json:"confirmation_sent_at,omitempty"`
	ConfirmedAt        *time.Time `json:"confirmed_at,omitempty"`

	// support additional user verification (beyond email)
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	RecoveryToken  string     `json:"-" gorm:"type:varchar(255)"`
	RecoverySentAt *time.Time `json:"recovery_sent_at,omitempty"`

	EmailChangeToken  string     `json:"-" gorm:"type:varchar(255)"`
	EmailChange       string     `json:"new_email,omitempty" gorm:"type:varchar(320)"`
	EmailChangeSentAt *time.Time `json:"email_change_sent_at,omitempty"`

	InvitedAt *time.Time `json:"invited_at,omitempty"`

	AppMetaData  Map `json:"app_metadata"`
	UserMetaData Map `json:"user_metadata"`

	IsSuperAdmin bool `json:"-"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// NewUser initializes a new user from an email, password and user data.
func NewUser(email, password string, userData map[string]interface{}) (*User, error) {
	id := uuid.New()
	pw, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	username := ""
	// TODO: make username a 1st class param
	if userData != nil {
		if val, has := userData["username"]; has {
			username = val.(string)
			delete(userData, "username")
		}
	}

	u := &User{
		ID:                id,
		Username:          username,
		Email:             email,
		UserMetaData:      userData,
		EncryptedPassword: pw,
	}
	return u, nil
}

func NewSystemUser() *User {
	return &User{
		ID:           SystemUserUUID,
		IsSuperAdmin: true,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.BeforeUpdate(tx)
}

func (u *User) BeforeUpdate(*gorm.DB) error {
	if u.ID == SystemUserUUID {
		return errors.New("cannot save system user")
	}
	return nil
}

func (u *User) BeforeSave(*gorm.DB) error {
	if u.ID == SystemUserUUID {
		return errors.New("cannot save system user")
	}

	if u.ConfirmedAt != nil && u.ConfirmedAt.IsZero() {
		u.ConfirmedAt = nil
	}
	if u.InvitedAt != nil && u.InvitedAt.IsZero() {
		u.InvitedAt = nil
	}
	if u.ConfirmationSentAt != nil && u.ConfirmationSentAt.IsZero() {
		u.ConfirmationSentAt = nil
	}
	if u.RecoverySentAt != nil && u.RecoverySentAt.IsZero() {
		u.RecoverySentAt = nil
	}
	if u.EmailChangeSentAt != nil && u.EmailChangeSentAt.IsZero() {
		u.EmailChangeSentAt = nil
	}
	if u.LastSignInAt != nil && u.LastSignInAt.IsZero() {
		u.LastSignInAt = nil
	}
	return nil
}

// SetRole sets the users Role to roleName
func (u *User) SetRole(tx *storage.Connection, roleName string) error {
	u.Role = strings.TrimSpace(roleName)
	return tx.Model(&u).Select("role").Updates(u).Error
}

// HasRole returns true when the users role is set to roleName
func (u *User) HasRole(roleName string) bool {
	return u.Role == roleName
}

// UpdateUserMetaData sets all user data from a map of updates,
// ensuring that it doesn't override attributes that are not
// in the provided map.
func (u *User) UpdateUserMetaData(tx *storage.Connection, updates map[string]interface{}) error {
	if u.UserMetaData == nil {
		u.UserMetaData = updates
	} else if updates != nil {
		for key, value := range updates {
			if value != nil {
				u.UserMetaData[key] = value
			} else {
				delete(u.UserMetaData, key)
			}
		}
	}
	return tx.Model(&u).Select("user_meta_data").Updates(u).Error
}

// UpdateAppMetaData updates all app data from a map of updates
func (u *User) UpdateAppMetaData(tx *storage.Connection, updates map[string]interface{}) error {
	if u.AppMetaData == nil {
		u.AppMetaData = updates
	} else if updates != nil {
		for key, value := range updates {
			if value != nil {
				u.AppMetaData[key] = value
			} else {
				delete(u.AppMetaData, key)
			}
		}
	}
	return tx.Model(&u).Select("app_meta_data").Updates(u).Error
}

func (u *User) SetEmail(tx *storage.Connection, email string) error {
	u.Email = email
	return tx.Model(&u).Select("email").Updates(u).Error
}

// hashPassword generates a hashed password from a plaintext string
func hashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pw), nil
}

func (u *User) UpdatePassword(tx *storage.Connection, password string) error {
	pw, err := hashPassword(password)
	if err != nil {
		return err
	}
	u.EncryptedPassword = pw
	// note that the user has changed their password
	u.ChangePassword = false
	return tx.Model(&u).Select("change_password", "encrypted_password").Updates(u).Error
}

// Authenticate a user from a password
func (u *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	return err == nil
}

// Confirm resets the email confirmation token and the confirm timestamp
func (u *User) Confirm(tx *storage.Connection) error {
	u.ConfirmationToken = ""
	now := time.Now()
	u.ConfirmedAt = &now
	return tx.Model(&u).Select("confirmation_token", "confirmed_at").Updates(u).Error
}

// ConfirmEmailChange confirm the change of email for a user
func (u *User) ConfirmEmailChange(tx *storage.Connection) error {
	u.Email = u.EmailChange
	u.EmailChange = ""
	u.EmailChangeToken = ""
	return tx.Model(&u).Select("email", "email_change", "email_change_token").Updates(u).Error
}

// IsConfirmed checks if a user has confirmed their email address.
func (u *User) IsConfirmed() bool {
	return u.ConfirmedAt != nil
}

// Verify sets the verification timestamp
func (u *User) Verify(tx *storage.Connection) error {
	now := time.Now()
	u.VerifiedAt = &now
	return tx.Model(&u).Select("verified_at").Updates(u).Error
}

// IsVerified checks if a user has been verified
func (u *User) IsVerified() bool {
	return u.VerifiedAt != nil
}

// Recover resets the recovery token
func (u *User) Recover(tx *storage.Connection) error {
	u.RecoveryToken = ""
	// make the user also have to set a new password now
	u.ChangePassword = true
	return tx.Model(&u).Select("change_password", "recovery_token").Updates(u).Error
}

// Deprecated: findUser is going away. we should build the query as part of the tx
func findUser(tx *storage.Connection, query string, args ...interface{}) (*User, error) {
	obj := &User{}
	if err := tx.Where(query, args...).First(obj).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, UserNotFoundError{}
		}
		err = fmt.Errorf("%w finding user", err)
		return nil, err
	}

	return obj, nil
}

// FindUserByConfirmationToken finds users with the matching confirmation token.
func FindUserByConfirmationToken(tx *storage.Connection, token string) (*User, error) {
	return findUser(tx, "confirmation_token = ?", token)
}

// FindUserByEmail finds a user with the matching email.
func FindUserByEmail(tx *storage.Connection, email string) (*User, error) {
	return findUser(tx, "email = ?", email)
}

// FindUserByID finds a user matching the provided ID.
func FindUserByID(tx *storage.Connection, id uuid.UUID) (*User, error) {
	return findUser(tx, "id = ?", id)
}

// FindUserByRecoveryToken finds a user with the matching recovery token.
func FindUserByRecoveryToken(tx *storage.Connection, token string) (*User, error) {
	return findUser(tx, "recovery_token = ?", token)
}

// FindUserWithRefreshToken finds a user from the provided refresh token.
func FindUserWithRefreshToken(tx *storage.Connection, token string) (*User, *RefreshToken, error) {
	refreshToken := &RefreshToken{}
	if err := tx.First(refreshToken, "token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, RefreshTokenNotFoundError{}
		}
		err = fmt.Errorf("%w finding refresh token", err)
		return nil, nil, err
	}

	u, err := findUser(tx, "id = ?", refreshToken.UserID)
	if err != nil {
		return nil, nil, err
	}

	return u, refreshToken, nil
}

// FindUsers finds users
func FindUsers(tx *storage.Connection, pageParams *Pagination, sortParams *SortParams, filter string) ([]*User, error) {
	var users []*User
	q := tx.Model(users)

	if filter != "" {
		lf := "%" + filter + "%"
		q = q.Where("(email LIKE ? OR user_meta_data COLLATE utf8mb4_unicode_ci LIKE ?)",
			lf, lf)
	}

	if sortParams != nil && len(sortParams.Fields) > 0 {
		for _, field := range sortParams.Fields {
			q = q.Order(field.Name + " " + string(field.Dir))
		}
	}

	var err error
	if pageParams != nil {
		p := paginator.New(adapter.NewGORMAdapter(q), int(pageParams.PerPage))
		p.SetPage(int(pageParams.Page))
		if err = p.Results(&users); err != nil {
			return nil, err
		}
		var cnt int
		if cnt, err = p.PageNums(); err != nil {
			return nil, err
		}
		pageParams.Count = uint64(cnt)
	} else {
		err = q.Find(&users).Error
	}

	return users, err
}

// IsDuplicatedEmail returns whether a user exists with a matching email.
func IsDuplicatedEmail(tx *storage.Connection, email string) (bool, error) {
	_, err := FindUserByEmail(tx, email)
	if err != nil {
		if IsNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// FindUserByUsername finds a user with the matching username.
func FindUserByUsername(tx *storage.Connection, name string) (*User, error) {
	return findUser(tx, "username = ?", name)
}

// IsDuplicatedUsername returns whether a user exists with a matching username.
func IsDuplicatedUsername(tx *storage.Connection, name string) (bool, error) {
	_, err := FindUserByUsername(tx, name)
	if err != nil {
		if IsNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

/*
// IsDuplicatedEmail returns whether a user exists with a matching email.
func IsDuplicatedEmail(tx *storage.Connection, email string) (bool, error) {
	u, err := FindUserByEmail(tx, email)
	if errors.Is(err, storage.ErrNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return u != nil, nil
}

*/
