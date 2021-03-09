package migration

import (
	"testing"

	"github.com/jrapoport/gothic/test/tdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlan_Run(t *testing.T) {
	db := tdb.DB(t)
	ma := ModelA{}
	mb := ModelB{}
	migA := NewMigration("A", ma)
	require.NotNil(t, migA)
	migB := NewMigration("B", mb)
	require.NotNil(t, migB)
	p := Plan{}
	err := p.Run(db, false)
	assert.NoError(t, err)
	p.AddMigration(migA)
	p.AddMigration(migB)
	err = p.Run(db, true)
	assert.NoError(t, err)
	p.AddMigrations([]*Migration{migA, migB})
	err = p.Run(db, true)
	assert.Error(t, err)
}
