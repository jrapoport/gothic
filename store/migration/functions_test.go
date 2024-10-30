package migration

import (
	"errors"
	"regexp"
	"testing"

	"github.com/jrapoport/gothic/test/tdb"
	"github.com/stretchr/testify/assert"
)

func TestNewMigrateFunc(t *testing.T) {
	t.Parallel()
	db := tdb.DB(t)
	m := ModelA{}
	fn := NewMigrateFunc(m)
	err := fn(db)
	assert.NoError(t, err)
	has := db.Migrator().HasTable(m)
	assert.True(t, has)
}

func TestNewMigrateFunc_Error(t *testing.T) {
	t.Parallel()
	db, mock := tdb.MockDB(t)
	create := "CREATE TABLE `model_as` (`id` bigint unsigned AUTO_INCREMENT," +
		"`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL," +
		"`deleted_at` datetime(3) NULL,`value` longtext,PRIMARY KEY (`id`)," +
		"INDEX `idx_model_as_deleted_at` (`deleted_at`))"
	mock.ExpectExec(regexp.QuoteMeta(create)).
		WillReturnError(errors.New("mock error"))
	m := ModelA{}
	fn := NewMigrateFunc(m)
	err := fn(db)
	assert.Error(t, err)
}

func TestNewMigrateWithIndexes(t *testing.T) {
	t.Parallel()
	const tableIdx = "idx_model_bs_value"
	var indexes = []string{ModelBIndex}
	db := tdb.DB(t)
	m := ModelB{}
	fn := NewMigrateFuncWithIndexes(m, indexes)
	err := fn(db)
	assert.NoError(t, err)
	tests := []struct {
		idx       string
		assertHas assert.BoolAssertionFunc
	}{
		{ModelBIndex, assert.False},
		{"idx_nope", assert.False},
		{tableIdx, assert.True},
	}
	for _, test := range tests {
		has := db.Migrator().HasIndex(m, test.idx)
		test.assertHas(t, has)
		has = db.Migrator().HasIndex(m, tableIdx)
		assert.True(t, has)
	}
	for _, test := range tests {
		fn = NewMigrateFuncWithIndexes(m, []string{test.idx})
		err = fn(db)
		has := db.Migrator().HasIndex(m, test.idx)
		test.assertHas(t, has)
		has = db.Migrator().HasIndex(m, tableIdx)
		assert.True(t, has)
	}
	err = migrateIndexes(db, nil, indexes)
	assert.Error(t, err)
}

func TestNewRollbackFunc(t *testing.T) {
	t.Parallel()
	db := tdb.DB(t)
	m := &ModelA{}
	err := db.AutoMigrate(m)
	assert.NoError(t, err)
	has := db.Migrator().HasTable(m)
	assert.True(t, has)
	rb := NewRollbackFunc(m)
	err = rb(db)
	assert.NoError(t, err)
	has = db.Migrator().HasTable(m)
	assert.False(t, has)
	// table name error
	rb = NewRollbackFunc(nil)
	err = rb(db)
	assert.Error(t, err)
}

/*
	fetchIndexes := func(tx *gorm.DB) []string {
		var indexNames []string
		query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type = 'index' AND tbl_name = '%s'", tableName)
		if err := tx.Raw(query).Scan(&indexNames).Error; err != nil {
			panic(err)
		}
		return indexNames
	}
*/
