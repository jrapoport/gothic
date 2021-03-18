package core

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/codes"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/mail"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

// OpenMail opens the mail client.
func (a *API) OpenMail() error {
	var err error
	if a.mail != nil {
		a.CloseMail()
	}
	a.mail, err = mail.NewMailClient(a.config, a.log)
	if err != nil {
		return a.logError(err)
	}
	if a.mail.IsOffline() {
		a.log.Warn("mail is offline")
		return nil
	}
	err = a.mail.Open()
	if err != nil {
		return a.logError(err)
	}
	a.log.Info("mail open")
	return nil
}

// CloseMail closes the mail client.
func (a *API) CloseMail() {
	defer func() {
		a.mail = nil
	}()
	if a.mail == nil {
		return
	}
	err := a.mail.Close()
	if err != nil {
		err = fmt.Errorf("failed to close mail %s", err)
		a.log.Error(err)
		return
	}
	a.log.Info("mail closed")
}

// ValidateEmail returns nil if the email is valid,
// otherwise an error indicating the reason it is invalid
func (a *API) ValidateEmail(email string) (string, error) {
	var s string
	var err error
	if a.mail == nil || a.mail.IsOffline() {
		s, err = validate.Email(email)
	} else {
		s, err = a.mail.ValidateEmail(email)
	}
	return s, a.logError(err)
}

// SendConfirmUser sends a invite mail to a new user
func (a *API) SendConfirmUser(ctx context.Context, userID uuid.UUID) error {
	if a.mail == nil || a.mail.IsOffline() {
		return nil
	}
	return a.sendConfirmToken(ctx, userID,
		func(u *user.User) (bool, error) {
			if u.IsConfirmed() {
				return false, nil
			}
			if !u.IsRestricted() {
				err := errors.New("invalid user")
				return false, err
			}
			return true, nil
		},
		func(to string, ct *token.ConfirmToken) error {
			referrerURL := a.config.Mail.ConfirmUser.ReferralURL
			return a.mail.SendConfirmUser(to, ct.String(), referrerURL)
		})
}

// SendResetPassword sends a password recovery mail
func (a *API) SendResetPassword(ctx context.Context, userID uuid.UUID) error {
	if a.mail == nil || a.mail.IsOffline() {
		return nil
	}
	return a.sendConfirmToken(ctx, userID,
		func(u *user.User) (bool, error) {
			if u.IsLocked() {
				err := errors.New("invalid user")
				return false, err
			}
			return true, nil
		},
		func(to string, ct *token.ConfirmToken) error {
			referrerURL := a.config.Mail.ResetPassword.ReferralURL
			return a.mail.SendResetPassword(to, ct.String(), referrerURL)
		})
}

// SendChangeEmail sends an email change confirmation mail to a user
func (a *API) SendChangeEmail(ctx context.Context, userID uuid.UUID, newAddress string) error {
	if a.mail == nil || a.mail.IsOffline() {
		return nil
	}
	return a.sendConfirmToken(ctx, userID,
		func(u *user.User) (bool, error) {
			if !u.IsActive() {
				err := errors.New("invalid user")
				return false, err
			}
			return true, nil
		},
		func(to string, ct *token.ConfirmToken) error {
			tok, err := jwt.NewSignedData(a.config.JWT, types.Map{
				key.Token: ct.String(),
				key.Email: newAddress,
			})
			if err != nil {
				return err
			}
			referrerURL := a.config.Mail.ChangeEmail.ReferralURL
			return a.mail.SendChangeEmail(to, newAddress, tok, referrerURL)
		})
}

// SendInviteUser sends a invite mail to a user
func (a *API) SendInviteUser(ctx context.Context, userID uuid.UUID, toAddress string) error {
	if a.mail == nil || a.mail.IsOffline() {
		return nil
	}
	if userID != user.SystemID && !a.config.Signup.CanSendInvites() {
		err := errors.New("invites disabled")
		return err
	}
	err := a.sendSignupCode(ctx, userID, toAddress,
		func(tx *store.Connection) (*code.SignupCode, error) {
			return codes.CreateSignupCode(tx, userID, code.Invite, 1, true)
		},
		func(from, to string, sc *code.SignupCode) error {
			tok, err := jwt.NewSignedData(a.config.JWT, types.Map{
				key.Token: sc.Code(),
				key.Email: to,
			})
			if err != nil {
				return err
			}
			referrerURL := a.config.Mail.InviteUser.ReferralURL
			return a.mail.SendInviteUser(from, to, tok, referrerURL)
		})
	return a.logError(err)
}
