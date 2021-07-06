package audit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/config/provider"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/account"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"github.com/jrapoport/gothic/store/types/key"
)

func TestLogSignup(t *testing.T) {
	r := user.RoleAdmin
	testLogEntry(t, auditlog.Signup, uuid.New(),
		types.Map{
			key.Role: r.String(),
		},
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogSignup(ctx, conn, uid, r)
		})
}

func TestLogCodeSent(t *testing.T) {
	uid := uuid.New()
	sc := code.NewSignupCode(uid, code.PIN, 1)
	sc.ID = 100
	sc.CreatedAt = time.Now().UTC()
	testLogEntry(t, auditlog.CodeSent, uid, logToken(sc),
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogCodeSent(ctx, conn, sc)
		})
}

func TestLogConfirmationSent(t *testing.T) {
	uid := uuid.New()
	tk := token.NewConfirmToken(uid, time.Second)
	tk.ID = 100
	tk.CreatedAt = time.Now().UTC()
	testLogEntry(t, auditlog.ConfirmSent, uid, logToken(tk),
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogConfirmSent(ctx, conn, tk)
		})
}

func TestLogConfirmed(t *testing.T) {
	testLogEntry(t, auditlog.Confirmed, uuid.New(), nil,
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogConfirmed(ctx, conn, uid)
		})
}

func TestLogBanned(t *testing.T) {
	testLogEntry(t, auditlog.Banned, uuid.New(), nil,
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogBanned(ctx, conn, uid)
		})
}

func TestLogDeleted(t *testing.T) {
	testLogEntry(t, auditlog.Deleted, uuid.New(), nil,
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogDeleted(ctx, conn, uid)
		})
}

func TestLogLinked(t *testing.T) {
	la := &account.Linked{
		Type:      account.Auth,
		Provider:  provider.Google,
		AccountID: uuid.New().String(),
		Data: types.Map{
			key.IPAddress: testIPAddress,
		},
	}
	fields := types.Map{
		key.Provider:  provider.Google,
		key.AccountID: la.AccountID,
		key.Type:      la.Type.String(),
	}
	testLogEntry(t, auditlog.Linked, uuid.New(), fields,
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, fields types.Map) error {
			return LogLinked(ctx, conn, uid, la)
		})
}
