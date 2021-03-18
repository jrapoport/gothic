package audit

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
)

// SearchEntries search the audit log entries.
func SearchEntries(conn *store.Connection, s store.Sort, f store.Filters, p *store.Pagination) ([]*auditlog.AuditLog, error) {
	tx := conn.Model(auditlog.AuditLog{})
	if uid, ok := f[key.UserID].(string); ok && uid != "" {
		id, err := uuid.Parse(uid)
		if err != nil {
			return nil, err
		}
		f[key.UserID] = id
	}
	if typ, ok := f[key.Type].(string); ok && typ != "" {
		f[key.Type] = auditlog.TypeFromString(typ)
	}
	flt := store.Filter{
		Filters:   f,
		DataField: key.Fields,
		Fields: []string{
			key.Action,
			key.Type,
			key.UserID,
		},
	}
	var logs []*auditlog.AuditLog
	err := store.Search(tx, &logs, s, flt, p)
	if err != nil {
		return nil, err
	}
	return logs, err
}
