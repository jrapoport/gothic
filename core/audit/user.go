package audit

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

// LogLogin log user login
func LogLogin(ctx context.Context, conn *store.Connection, userID uuid.UUID) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Login, userID, nil)
	return err
}

// LogLogout log user logout
func LogLogout(ctx context.Context, conn *store.Connection, userID uuid.UUID) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Logout, userID, nil)
	return err
}

// LogPasswordChange log user password change
func LogPasswordChange(ctx context.Context, conn *store.Connection, userID uuid.UUID) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Password, userID, nil)
	return err
}

// LogEmailChange log user email change
func LogEmailChange(ctx context.Context, conn *store.Connection, userID uuid.UUID) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Email, userID, nil)
	return err
}

// LogUserUpdated log user updated
func LogUserUpdated(ctx context.Context, conn *store.Connection, userID uuid.UUID) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Updated, userID, nil)
	return err
}

// LogChangeRole log user role change
func LogChangeRole(ctx context.Context, conn *store.Connection, userID uuid.UUID, r user.Role) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.ChangeRole, userID, types.Map{
		key.Role: r.String(),
	})
	return err
}
