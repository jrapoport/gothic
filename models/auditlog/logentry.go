package auditlog

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"gorm.io/gorm"
)

func init() {
	store.AddAutoMigration("1000-audit_logs", AuditLog{})
}

// AuditLog is the database model for audit log entries.
type AuditLog struct {
	gorm.Model
	Type   Type      `json:"type"`
	Action Action    `json:"action"`
	UserID uuid.UUID `json:"user_id" gorm:"index;type:char(36)"`
	Fields types.Map `json:"fields"`
}

// NewAuditLog returns a new log entry
func NewAuditLog(t Type, a Action, userID uuid.UUID, fields types.Map) *AuditLog {
	e := &AuditLog{
		Type:   t,
		Action: a,
		UserID: userID,
		Fields: fields,
	}
	return e
}
