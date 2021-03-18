package audit

import (
	"time"

	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
)

// LogTokenGranted log token granted
func LogTokenGranted(ctx context.Context, conn *store.Connection, t token.Token) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Granted, t.IssuedTo(), logToken(t))
	return err
}

// LogTokenRefreshed log token refreshed
func LogTokenRefreshed(ctx context.Context, conn *store.Connection, t token.Token) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Refreshed, t.IssuedTo(), logToken(t))
	return err
}

// LogTokenRevoked log token revoked
func LogTokenRevoked(ctx context.Context, conn *store.Connection, t token.Token) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.Revoked, t.IssuedTo(), logToken(t))
	return err
}

/*
// LogRevokedAll log all tokens revoked
func LogRevokedAll(ctx context.Context, conn *store.Connection, t token.Token) error {
	_, err := CreateLogEntry(ctx, conn, auditlog.RevokedAll, t.IssuedTo(), logToken(t))
	return err
}
*/

func logToken(t token.Token) types.Map {
	data := types.Map{
		key.Token:  t.String(),
		key.Class:  t.Class().String(),
		key.Usage:  t.Usage().String(),
		key.Issued: t.Issued().UTC().Format(time.RFC3339),
	}
	if !t.LastUsed().IsZero() {
		data[key.LastUsed] = t.LastUsed().UTC().Format(time.RFC3339)
	}
	if !t.ExpirationDate().IsZero() {
		data[key.ExpirationDate] = t.ExpirationDate().UTC().Format(time.RFC3339)
	}
	if !t.Revoked().IsZero() {
		data[key.Revoked] = t.Revoked().UTC().Format(time.RFC3339)
	}
	return data
}
