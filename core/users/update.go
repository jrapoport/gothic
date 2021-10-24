package users

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/imdario/mergo"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

// Update updates a user
func Update(conn *store.Connection, u *user.User, username *string, data types.Map) (bool, error) {
	if u == nil || !u.IsActive() {
		return false, errors.New("invalid user")
	}
	if username != nil {
		u.Username = *username
	} else if data == nil || len(data) <= 0 {
		return false, nil
	}
	err := mergo.Merge(&u.Data, data, mergo.WithOverride)
	if err != nil {
		return false, err
	}
	err = conn.Model(u).Select(key.Username, key.Data).Updates(u).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateMetadata updates a user's metadata
func UpdateMetadata(conn *store.Connection, u *user.User, meta types.Map) (bool, error) {
	// don't check if the user is active since we should
	// be able to operate on banned users, etc.
	if u == nil {
		return false, errors.New("invalid user")
	}
	if meta == nil || len(meta) <= 0 {
		return false, nil
	}
	for k := range meta {
		if _, ok := key.Reserved[k]; ok {
			delete(meta, k)
		}
	}
	err := mergo.Merge(&u.Metadata, meta, mergo.WithOverride)
	if err != nil {
		return false, err
	}
	err = conn.Model(u).Select(key.Metadata).Updates(u).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConfirmUser confirms a user.
func ConfirmUser(conn *store.Connection, u *user.User, t time.Time) error {
	if u == nil || u.IsLocked() {
		return errors.New("invalid user")
	}
	if u.IsConfirmed() {
		return nil
	}
	if t.IsZero() {
		t = time.Now().UTC()
	}
	u.ConfirmedAt = &t
	u.Status = user.Active
	return conn.Model(u).Select(key.ConfirmedAt, key.Status).Updates(u).Error
}

// ConfirmIfNeeded confirms a user, returns true if the user was previously unconfirmed
func ConfirmIfNeeded(conn *store.Connection, ct *token.ConfirmToken, u *user.User) (bool, error) {
	if ct == nil || !ct.Usable() {
		return false, errors.New("invalid token")
	}
	if u == nil || u.IsLocked() {
		return false, errors.New("invalid user")
	}
	var confirmed bool
	err := conn.Transaction(func(tx *store.Connection) error {
		err := tokens.UseToken(tx, ct)
		if err != nil {
			return err
		}
		if u.IsConfirmed() {
			return nil
		}
		err = ConfirmUser(tx, u, *ct.UsedAt)
		if err != nil {
			return err
		}
		confirmed = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return confirmed, nil
}

// ChangeRole changes the role for a user.
func ChangeRole(conn *store.Connection, u *user.User, r user.Role) error {
	if u == nil || u.IsLocked() {
		return errors.New("invalid user")
	}
	if !r.Valid() {
		return errors.New("invalid role")
	}
	if u.Role == r {
		return nil
	}
	u.Role = r
	return conn.Model(u).Update(key.Role, u.Role).Error
}

// ChangeEmail changes the email for a user.
func ChangeEmail(conn *store.Connection, u *user.User, email string) error {
	if u == nil || !u.IsActive() {
		return errors.New("invalid user")
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email: %s", email)
	}
	if u.Email == addr.Address {
		return nil
	}
	u.Email = addr.Address
	return conn.Model(&u).Update(key.Email, u.Email).Error
}

// ChangePassword changes the password for a user.
func ChangePassword(conn *store.Connection, u *user.User, pw string) error {
	if u == nil || u.IsLocked() {
		return errors.New("invalid user")
	}
	if u.Provider.IsExternal() {
		return errors.New("invalid provider")
	}
	hashed := utils.HashPassword(pw)
	u.Password = hashed
	return conn.Model(&u).Update(key.Password, u.Password).Error
}

// LockUser locks a user.
func LockUser(conn *store.Connection, u *user.User) error {
	// do not step on banned users
	if u == nil || u.IsBanned() {
		return nil
	}
	u.Status = user.Locked
	return conn.Model(u).Update(key.Status, u.Status).Error
}

// BanUser bans a user.
func BanUser(conn *store.Connection, u *user.User) error {
	// step on any other user state
	if u == nil || !u.Valid() {
		return nil
	}
	u.Status = user.Banned
	return conn.Model(u).Update(key.Status, u.Status).Error
}

// DeleteUser deletes a user.
func DeleteUser(conn *store.Connection, u *user.User, hard bool) error {
	if u == nil || !u.Valid() {
		return nil
	}
	// we don't want to delete a banned user so their data
	// stays banned and they can't reuse an email address
	if u.IsBanned() {
		return nil
	}
	db := conn.DB
	if hard {
		db = db.Unscoped()
	}
	return db.Delete(u).Error
}
