package store

import "github.com/jrapoport/gothic/store/migration"

var plan migration.Plan

// AddAutoMigrationToPlan adds a model to the global migration plan
func AddAutoMigrationToPlan(m *migration.Migration) {
	plan.AddMigration(m)
}

// AddAutoMigration adds a model to the global migration plan
func AddAutoMigration(id string, model interface{}) {
	m := migration.NewMigration(id, model)
	AddAutoMigrationToPlan(m)
}

// AddAutoMigrationWithIndexes adds a model to the global migration plan with indexes
// TODO: can we make the namespacing of indexes a default in gorm or the migration?
func AddAutoMigrationWithIndexes(id string, model interface{}, indexes []string) {
	m := migration.NewMigrationWithIndexes(id, model, indexes)
	AddAutoMigrationToPlan(m)
}
