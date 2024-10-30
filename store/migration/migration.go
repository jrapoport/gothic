package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migration represents a database migration (a modification to be made on the database).
type Migration struct {
	gormigrate.Migration
	model   interface{}
	indexes []string
}

// Run returns the migration on the db.
func (m *Migration) Run(db *gorm.DB) error {
	plan := NewPlan()
	plan.AddMigration(m)
	return plan.Run(db, false)
}

// NewMigration returns a new Migration for the model with id.
func NewMigration(id string, model interface{}) *Migration {
	return NewMigrationWithIndexes(id, model, nil)
}

// NewMigrationWithIndexes returns a new Migration for the indexed model with id.
func NewMigrationWithIndexes(id string, model interface{}, indexes []string) *Migration {
	mg := NewMigrateFunc(model)
	if len(indexes) > 0 {
		mg = NewMigrateFuncWithIndexes(model, indexes)
	}
	rb := NewRollbackFunc(model)
	return &Migration{
		Migration: gormigrate.Migration{
			ID:       id,
			Migrate:  gormigrate.MigrateFunc(mg),
			Rollback: gormigrate.RollbackFunc(rb),
		},
		model:   model,
		indexes: indexes,
	}
}
