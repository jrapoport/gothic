package core

import (
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/store"
)

// GetAuditLog returns the log entry for the id.
func (a *API) GetAuditLog(_ context.Context, id uint) (*auditlog.LogEntry, error) {
	return audit.GetLogEntry(a.conn, id)
}

// SearchAuditLogs searches the audit logs.
func (a *API) SearchAuditLogs(ctx context.Context, f store.Filters, page *store.Pagination) ([]*auditlog.LogEntry, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return audit.SearchEntries(a.conn, ctx.GetSort(), f, page)
}
