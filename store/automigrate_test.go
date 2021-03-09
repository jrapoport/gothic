package store

import (
	"testing"

	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAddAutoMigration(t *testing.T) {
	type ModA struct {
		gorm.Model
		Value string
	}
	type ModB struct {
		gorm.Model
		Value string `gorm:"index:idx_value"`
	}
	// ModelBIndex is the name of the index for ModelB
	const (
		BIndex   = "idx_value"
		tableIdx = "idx_test_mod_bs_value"
	)
	var ma ModA
	var mb ModB
	t.Cleanup(func() {
		plan.Clear()
	})
	AddAutoMigration("01", ma)
	AddAutoMigrationWithIndexes("02", mb, []string{BIndex})
	c := tconf.TempDB(t)
	conn, err := Dial(c, nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	err = conn.AutoMigrate()
	require.NoError(t, err)
	has := conn.DB.Migrator().HasTable(ma)
	assert.True(t, has)
	has = conn.DB.Migrator().HasTable(mb)
	assert.True(t, has)
	has = conn.DB.Migrator().HasIndex(mb, tableIdx)
	assert.True(t, has)
	err = conn.AutoMigrate()
	require.NoError(t, err)
}
