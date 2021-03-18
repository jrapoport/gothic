package migration

import (
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
}

func TestNewRollbackFunc(t *testing.T) {
	t.Parallel()
	db := tdb.DB(t)
	tc := &ModelA{}
	err := db.AutoMigrate(tc)
	assert.NoError(t, err)
	has := db.Migrator().HasTable(tc)
	assert.True(t, has)
	rb := NewRollbackFunc(tc)
	err = rb(db)
	assert.NoError(t, err)
	has = db.Migrator().HasTable(tc)
	assert.False(t, has)
}
