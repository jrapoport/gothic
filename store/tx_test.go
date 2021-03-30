package store

import (
	"math"
	"testing"

	"github.com/jrapoport/gothic/store/migration"
	"github.com/jrapoport/gothic/test/tdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHas(t *testing.T) {
	t.Parallel()
	type ModA struct {
		gorm.Model
		Value string
	}
	type ModB struct {
		gorm.Model
		Value string `gorm:"index:idx_value"`
	}
	const testValue = "test-value"
	db := tdb.DB(t)
	p := migration.NewPlan()
	p.AddMigrations([]*migration.Migration{
		migration.NewMigration("1", ModA{}),
		migration.NewMigration("2", ModB{}),
	})
	err := p.Run(db, false)
	require.NoError(t, err)
	var has bool
	ma := &ModA{
		Value: testValue,
	}
	err = db.Create(ma).Error
	assert.NoError(t, err)
	has, err = Has(db, ma)
	assert.NoError(t, err)
	assert.True(t, has)
	has, err = Has(db, ma, "Value = ?", testValue)
	assert.NoError(t, err)
	assert.True(t, has)
	err = db.Delete(ma).Error
	assert.NoError(t, err)
	has, err = Has(db, ma)
	assert.NoError(t, err)
	assert.False(t, has)
	// false
	ma.ID = math.MaxInt8 - 1
	has, err = Has(db, ma)
	assert.NoError(t, err)
	assert.False(t, has)
	// false
	mb := &ModB{}
	has, err = Has(db, mb)
	assert.NoError(t, err)
	assert.False(t, has)
	// error
	has, err = Has(db, nil)
	assert.Error(t, err)
	assert.False(t, has)
}

func TestHasLast(t *testing.T) {
	t.Parallel()
	type ModA struct {
		gorm.Model
		Value string
	}
	type ModB struct {
		gorm.Model
		Value string `gorm:"index:idx_value"`
	}
	const testValue = "test-value"
	db := tdb.DB(t)
	p := migration.NewPlan()
	p.AddMigrations([]*migration.Migration{
		migration.NewMigration("1", ModA{}),
		migration.NewMigration("2", ModB{}),
	})
	err := p.Run(db, false)
	require.NoError(t, err)
	var has bool
	ma := &ModA{
		Value: testValue,
	}
	err = db.Create(ma).Error
	assert.NoError(t, err)
	has, err = HasLast(db, ma)
	assert.NoError(t, err)
	assert.True(t, has)
	has, err = HasLast(db, ma, "Value = ?", testValue)
	assert.NoError(t, err)
	assert.True(t, has)
	err = db.Delete(ma).Error
	assert.NoError(t, err)
	has, err = HasLast(db, ma)
	assert.NoError(t, err)
	assert.False(t, has)
	// false
	ma.ID = math.MaxInt8 - 1
	has, err = HasLast(db, ma)
	assert.NoError(t, err)
	assert.False(t, has)
	// false
	mb := &ModB{}
	has, err = HasLast(db, mb)
	assert.NoError(t, err)
	assert.False(t, has)
	// error
	has, err = HasLast(db, nil)
	assert.Error(t, err)
	assert.False(t, has)
}
