package users

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/imdario/mergo"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/utils"
)

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
		err := token.UseToken(tx, ct)
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
	if u.Provider.IsExternal() {
		return errors.New("invalid provider")
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email: %s", email)
	}
	u.Email = addr.Address
	return conn.Model(&u).Select(key.Email).Updates(u).Error
}

// ChangePassword changes the password for a user.
func ChangePassword(conn *store.Connection, u *user.User, pw string) error {
	if u == nil || u.IsLocked() {
		return errors.New("invalid user")
	}
	if u.Provider.IsExternal() {
		return errors.New("invalid provider")
	}
	hashed, err := utils.HashPassword(pw)
	if err != nil {
		return err
	}
	u.Password = hashed
	return conn.Model(&u).Select(key.Password).Updates(u).Error
}

// Update updates a user
func Update(conn *store.Connection, u *user.User, username *string, data types.Map) error {
	if u == nil || !u.IsActive() {
		return errors.New("invalid user")
	}
	if username != nil {
		u.Username = *username
	} else if data == nil {
		return nil
	}
	err := mergo.Merge(&u.Data, data, mergo.WithOverride)
	if err != nil {
		return err
	}
	return conn.Model(u).Select(key.Username, key.Data).Updates(u).Error
}

// LockUser locks a user.
func LockUser(conn *store.Connection, u *user.User) error {
	// do not step on banned users
	if u == nil || u.IsBanned() {
		return nil
	}
	u.Status = user.Locked
	return conn.Model(u).Select(key.Status).Updates(u).Error
}

// BanUser bans a user.
func BanUser(conn *store.Connection, u *user.User) error {
	// step on any other user state
	if u == nil || !u.Valid() {
		return nil
	}
	u.Status = user.Banned
	return conn.Model(u).Select(key.Status).Updates(u).Error
}

// DeleteUser deletes a user.
func DeleteUser(conn *store.Connection, u *user.User) error {
	if u == nil || !u.Valid() {
		return nil
	}
	// we don't want to delete a banned user so their data
	// stays banned and they can't reuse an email address
	if u.IsBanned() {
		return nil
	}
	return conn.Delete(u).Error
}
