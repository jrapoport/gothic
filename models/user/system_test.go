package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
)

func TestNewSystemUser(t *testing.T) {
	t.Parallel()
	conn, c := tconn.TempConn(t)
	su := NewSystemUser()
	assert.True(t, su.IsSystemUser())
	assert.Equal(t, RoleSystem, su.Role)
	err := conn.Create(su).Error
	assert.Error(t, err)
	su.ID = uuid.New()
	su.Provider = c.Provider()
	err = conn.Create(su).Error
	assert.NoError(t, err)
}
