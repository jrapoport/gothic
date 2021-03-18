package audit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/store"
)

func TestLogTokenGranted(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	tk := token.NewRefreshToken(uid)
	tk.ID = 100
	tk.CreatedAt = time.Now().UTC()
	testLogEntry(t, auditlog.Granted, uid, logToken(tk),
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogTokenGranted(ctx, conn, tk)
		})
}

func TestLogTokenRefreshed(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	tk := token.NewRefreshToken(uid)
	tk.ID = 100
	tk.CreatedAt = time.Now().UTC()
	testLogEntry(t, auditlog.Refreshed, uid, logToken(tk),
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogTokenRefreshed(ctx, conn, tk)
		})
}

func TestLogTokenRevoked(t *testing.T) {
	t.Parallel()
	uid := uuid.New()
	tk := token.NewRefreshToken(uid)
	tm := time.Now().UTC()
	tk = token.NewRefreshToken(uid)
	tk.ID = 100
	tk.DeletedAt.Time = tm
	tk.DeletedAt.Valid = true
	tk.ExpiredAt = &tm
	testLogEntry(t, auditlog.Revoked, uid, logToken(tk),
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogTokenRevoked(ctx, conn, tk)
		})
}

/*
func TestLogRevokedAll(t *testing.T) {
	uid := uuid.New()
	tk := token.NewRefreshToken(uid)
	testLogEntry(t, auditlog.RevokedAll, uid, logToken(tk),
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogRevokedAll(ctx, conn, tk)
		})
}
*/
