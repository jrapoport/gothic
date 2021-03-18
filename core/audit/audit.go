package audit

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
)

// CreateLogEntry creates a new audit log entry.
// TODO: support IP geolocation https://github.com/apilayer/freegeoip/
func CreateLogEntry(ctx context.Context, conn *store.Connection,
	action auditlog.Action, userID uuid.UUID, fields types.Map) (*auditlog.AuditLog, error) {
	t := action.Type()
	if ctx != nil {
		if fields == nil {
			fields = types.Map{}
		}
		if fields[key.IPAddress] == nil && ctx.GetIPAddress() != "" {
			fields[key.IPAddress] = ctx.GetIPAddress()
		}
		if fields[key.Provider] == nil &&
			ctx.GetProvider() != provider.Unknown {
			fields[key.Provider] = ctx.GetProvider()
		}
		if fields[key.AdminID] == nil &&
			ctx.GetAdminID() != uuid.Nil {
			fields[key.AdminID] = ctx.GetAdminID().String()
		} else if fields[key.UserID] == nil &&
			ctx.GetUserID() != uuid.Nil {
			fields[key.UserID] = ctx.GetUserID().String()
		}
	}
	// conn.Logger.LogMode(0).Info(ctx, "%s %s: %s: %v", t, action, userID, fields)
	le := auditlog.NewAuditLog(t, action, userID, fields)
	return le, conn.Create(le).Error
}

// GetLogEntry gets an audit log entry.
func GetLogEntry(conn *store.Connection, id uint) (*auditlog.AuditLog, error) {
	le := new(auditlog.AuditLog)
	le.ID = id
	err := conn.First(le).Error
	if err != nil {
		return nil, err
	}
	return le, nil
}
