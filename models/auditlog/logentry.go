package auditlog

import (
	"github.com/google/uuid"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/types"
	"gorm.io/gorm"
)

func init() {
	store.AddAutoMigration("1000-audit_log", LogEntry{})
}

// LogEntry is the database model for audit log entries.
type LogEntry struct {
	gorm.Model
	Type   Type      `json:"type"`
	Action Action    `json:"action"`
	UserID uuid.UUID `json:"user_id" gorm:"index;type:char(36)"`
	Fields types.Map `json:"fields"`
}

// NewLogEntry returns a new log entry
func NewLogEntry(t Type, a Action, userID uuid.UUID, fields types.Map) *LogEntry {
	e := &LogEntry{
		Type:   t,
		Action: a,
		UserID: userID,
		Fields: fields,
	}
	return e
}

// TableName overrides the table name
func (LogEntry) TableName() string {
	return "audit_log"
}
