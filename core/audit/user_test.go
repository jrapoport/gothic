package audit

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
)

func TestLogLogin(t *testing.T) {
	t.Parallel()
	testLogEntry(t, auditlog.Login, uuid.New(), nil,
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogLogin(ctx, conn, uid)
		})
}

func TestLogLogout(t *testing.T) {
	t.Parallel()
	testLogEntry(t, auditlog.Logout, uuid.New(), nil,
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogLogout(ctx, conn, uid)
		})
}

func TestLogPasswordChange(t *testing.T) {
	t.Parallel()
	testLogEntry(t, auditlog.Password, uuid.New(), nil,
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogPasswordChange(ctx, conn, uid)
		})
}

func TestLogEmailChange(t *testing.T) {
	t.Parallel()
	testLogEntry(t, auditlog.Email, uuid.New(), nil,
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogEmailChange(ctx, conn, uid)
		})
}

func TestLogUpdate(t *testing.T) {
	t.Parallel()
	testLogEntry(t, auditlog.Updated, uuid.New(), nil,
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogUserUpdated(ctx, conn, uid)
		})
}

func TestLogChangeRole(t *testing.T) {
	t.Parallel()
	r := user.RoleAdmin
	testLogEntry(t, auditlog.ChangeRole, uuid.New(),
		types.Map{
			key.Role: r.String(),
		},
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogChangeRole(ctx, conn, uid, r)
		})
}
