package migration

import (
	"testing"

	"github.com/jrapoport/gothic/test/tdb"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// ModelA is a model for tests.
type ModelA struct {
	gorm.Model
	Value string
}

// ModelB is a model with indexes for tests.
type ModelB struct {
	gorm.Model
	Value string `gorm:"index:idx_value"`
}

// ModelBIndex is the name of the index for ModelB
const ModelBIndex = "idx_value"

func TestMigration_Run(t *testing.T) {
	t.Parallel()
	const tableIdx = "idx_model_bs_value"
	var indexes = []string{ModelBIndex}
	db := tdb.DB(t)
	ma := ModelA{}
	mb := ModelB{}
	migA := NewMigration("A", ma)
	assert.NotNil(t, migA)
	err := migA.Run(db)
	assert.NoError(t, err)
	has := db.Migrator().HasTable(ma)
	assert.True(t, has)
	migB := NewMigrationWithIndexes("B", mb, indexes)
	assert.NotNil(t, migB)
	err = migB.Run(db)
	assert.NoError(t, err)
	has = db.Migrator().HasTable(mb)
	assert.True(t, has)
	has = db.Migrator().HasIndex(mb, tableIdx)
	assert.True(t, has)
	has = db.Migrator().HasIndex(mb, ModelBIndex)
	assert.False(t, has)
}
