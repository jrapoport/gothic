package migration

import (
	"testing"

	"github.com/jrapoport/gothic/test/tdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type ModelC struct {
	gorm.Model
	FK []ModelD
}

type ModelD struct {
	gorm.Model
	Number   string
	ModelCID uint
}

func TestPlan_Run(t *testing.T) {
	t.Parallel()
	db := tdb.DB(t)
	ma := ModelA{}
	mb := ModelB{}
	mc := ModelC{}
	md := ModelD{}
	migA := NewMigration("A", ma)
	require.NotNil(t, migA)
	migB := NewMigration("B", mb)
	require.NotNil(t, migB)
	migC := NewMigration("C", mc)
	require.NotNil(t, migA)
	migD := NewMigration("D", md)
	require.NotNil(t, migB)
	p := NewPlan()
	err := p.Run(db, false)
	assert.NoError(t, err)
	p.AddMigration(migA)
	p.AddMigration(migB)
	p.AddMigration(migC)
	p.AddMigration(migD)
	err = p.Run(db, true)
	assert.NoError(t, err)
	p.AddMigrations([]*Migration{migA, migB})
	err = p.Run(db, true)
	assert.Error(t, err)
	p.Clear()
	assert.Len(t, p.migs, 0)
}

func TestPlan_Run_Errors(t *testing.T) {
	db := tdb.DB(t)
	ma := ModelA{}
	migA := NewMigration("A", ma)
	require.NotNil(t, migA)
	p := NewPlan()
	err := p.Run(db, false)
	assert.NoError(t, err)
	p.AddMigration(migA)
	p.migs[0].model = nil
	err = p.Run(db, false)
	assert.Error(t, err)
}
