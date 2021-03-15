package audit

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
)

// LogSignup logs a user signup.
func LogSignup(ctx context.Context, conn *store.Connection, userID uuid.UUID, r user.Role) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Signup, userID, types.Map{
		key.Role: r.String(),
	})
	return err
}

// LogCodeSent logs a sent signup code.
func LogCodeSent(ctx context.Context, conn *store.Connection, t token.Token) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.CodeSent, t.IssuedTo(), logToken(t))
	return err
}

// LogConfirmSent logs a sent confirmation token.
func LogConfirmSent(ctx context.Context, conn *store.Connection, t token.Token) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.ConfirmSent, t.IssuedTo(), logToken(t))
	return err
}

// LogConfirmed logs a confirmed user.
func LogConfirmed(ctx context.Context, conn *store.Connection, userID uuid.UUID) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Confirmed, userID, nil)
	return err
}

// LogBanned logs a banned user.
func LogBanned(ctx context.Context, conn *store.Connection, userID uuid.UUID) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Banned, userID, nil)
	return err
}

// LogDeleted logs a deleted user.
func LogDeleted(ctx context.Context, conn *store.Connection, userID uuid.UUID) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Deleted, userID, nil)
	return err
}

// LogLinked logs a linked account.
func LogLinked(ctx context.Context, conn *store.Connection, userID uuid.UUID, la *account.LinkedAccount) error {
	data := types.Map{
		key.Type:      la.Type.String(),
		key.Provider:  la.Provider,
		key.AccountID: la.AccountID,
	}
	_, err := CreateLogEntry(ctx, conn, auditlog.Linked, userID, data)
	return err
}
