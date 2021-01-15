package api

import (
	"context"
	"fmt"
	"time"

	"github.com/jrapoport/gothic/crypto"
	"github.com/jrapoport/gothic/mailer"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
)

func sendConfirmation(tx *storage.Connection, u *models.User, mailer mailer.Mailer, maxFrequency time.Duration, referrerURL string) error {
	if u.ConfirmationSentAt != nil && !u.ConfirmationSentAt.Add(maxFrequency).Before(time.Now()) {
		return nil
	}

	oldToken := u.ConfirmationToken
	u.ConfirmationToken = crypto.SecureToken()
	now := time.Now()
	if err := mailer.ConfirmationMail(u, referrerURL); err != nil {
		u.ConfirmationToken = oldToken
		err = fmt.Errorf("sending confirmation email %w", err)
		return err
	}
	u.ConfirmationSentAt = &now
	err := tx.Model(&u).Select("confirmation_token", "confirmation_sent_at").Updates(u).Error
	if err != nil {
		err = fmt.Errorf("updating user for confirmation %w", err)
	}
	return err
}

func sendInvite(tx *storage.Connection, u *models.User, mailer mailer.Mailer, referrerURL string) error {
	oldToken := u.ConfirmationToken
	u.ConfirmationToken = crypto.SecureToken()
	now := time.Now()
	if err := mailer.InviteMail(u, referrerURL); err != nil {
		u.ConfirmationToken = oldToken
		err = fmt.Errorf("sending invite email %w", err)
		return err
	}
	u.InvitedAt = &now
	err := tx.Model(&u).Select("confirmation_token", "invited_at").Updates(u).Error
	if err != nil {
		err = fmt.Errorf("updating user for invite %w", err)
	}
	return err
}

func (a *API) sendPasswordRecovery(tx *storage.Connection, u *models.User, mailer mailer.Mailer, maxFrequency time.Duration, referrerURL string) error {
	if u.RecoverySentAt != nil && !u.RecoverySentAt.Add(maxFrequency).Before(time.Now()) {
		return nil
	}

	oldToken := u.RecoveryToken
	u.RecoveryToken = crypto.SecureToken()
	now := time.Now()
	if err := mailer.RecoveryMail(u, referrerURL); err != nil {
		u.RecoveryToken = oldToken
		err = fmt.Errorf("sending recovery email %w", err)
		return err
	}
	u.RecoverySentAt = &now
	err := tx.Model(&u).Select("recovery_token", "recovery_sent_at").Updates(u).Error
	if err != nil {
		err = fmt.Errorf("updating user for recovery %w", err)
	}
	return err
}

func (a *API) sendEmailChange(tx *storage.Connection, u *models.User, mailer mailer.Mailer, email string, referrerURL string) error {
	oldToken := u.EmailChangeToken
	oldEmail := u.EmailChange
	u.EmailChangeToken = crypto.SecureToken()
	u.EmailChange = email
	now := time.Now()
	if err := mailer.EmailChangeMail(u, referrerURL); err != nil {
		u.EmailChangeToken = oldToken
		u.EmailChange = oldEmail
		return err
	}
	u.EmailChangeSentAt = &now
	err := tx.Model(&u).Select("email_change_token", "email_change", "email_change_sent_at").Updates(u).Error
	if err != nil {
		err = fmt.Errorf("updating user for email change %w", err)
	}
	return err
}

func (a *API) validateEmail(ctx context.Context, email string) error {
	if email == "" {
		return unprocessableEntityError("An email address is required")
	}
	m := a.Mailer(ctx)
	err := m.ValidateEmail(email)
	if err != nil {
		return unprocessableEntityError("Unable to validate email address: " + err.Error())
	}
	return nil
}
