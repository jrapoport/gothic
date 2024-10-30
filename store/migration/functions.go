package migration

import (
	"strings"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jrapoport/gothic/store/drivers"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MigrateFunc is the func signature for migrating.
type MigrateFunc gormigrate.MigrateFunc

// RollbackFunc is the func signature for rollbacking.
type RollbackFunc gormigrate.RollbackFunc

// NewMigrateFunc returns a new MigrateFunc for migrating.
func NewMigrateFunc(dst interface{}) MigrateFunc {
	return NewMigrateFuncWithIndexes(dst, nil)
}

// NewMigrateFuncWithIndexes returns a new MigrateFunc for migrating with namespaces indexes.
func NewMigrateFuncWithIndexes(dst interface{}, indexes []string) MigrateFunc {
	return func(tx *gorm.DB) error {
		err := tx.AutoMigrate(dst)
		if err != nil {
			return err
		}
		if indexes == nil {
			return nil
		}
		return migrateIndexes(tx, dst, indexes)
	}
}

// NewRollbackFunc returns a new RollbackFunc for rollbacking.
func NewRollbackFunc(src interface{}) RollbackFunc {
	return func(tx *gorm.DB) error {
		name, err := tableName(tx, src)
		if err != nil {
			return err
		}
		return tx.Migrator().DropTable(name)
	}
}

func migrateIndexes(tx *gorm.DB, dst interface{}, indexes []string) error {
	name, err := tableName(tx, dst)
	if err != nil {
		return err
	}
	const idxPrefix = "idx_"
	var nsPrefix = idxPrefix + name + "_"
	for _, idx := range indexes {
		if !tx.Migrator().HasIndex(dst, idx) {
			continue
		}
		if strings.HasPrefix(idx, nsPrefix) {
			continue
		}
		canonicalizeIndex := func(idx string) string {
			idx = strings.TrimPrefix(idx, idxPrefix)
			return nsPrefix + idx
		}
		newIdx := canonicalizeIndex(idx)
		if !tx.Migrator().HasIndex(dst, newIdx) {
			err = tx.Migrator().RenameIndex(dst, idx, newIdx)
			if err != nil {
				return err
			}
		}
		if tx.Name() == drivers.SQLite {
			err = dropIndexIfExists(tx, idx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dropIndexIfExists(tx *gorm.DB, name string) error {
	return tx.Exec("DROP INDEX IF EXISTS ?", clause.Column{Name: name}).Error
}

// tableName returns the table name for the model
func tableName(tx *gorm.DB, model interface{}) (string, error) {
	stmt := tx.Unscoped().Statement
	if err := stmt.Parse(model); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

/*
// NamespacedTable returns the table name with namespacing applied.
func NamespacedTable(tx *gorm.DB, dst interface{}) *gorm.DB {
	if mod, ok := dst.(schema.Tabler); ok {
		n := tx.NamingStrategy.TableName(mod.TableName())
		tx = tx.Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Table(n)
		})
	}
	return tx
}
*/
