package user

import (
	"testing"

	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
)

func TestNewSuperAdmin(t *testing.T) {
	const password = "test"
	conn, c := tconn.TempConn(t)
	sa := NewSuperAdmin(password)
	err := conn.Create(sa).Error
	assert.Error(t, err)
	sa.Provider = c.Provider()
	err = conn.Create(sa).Error
	assert.NoError(t, err)
	sa = NewSuperAdmin("password")
	err = conn.Create(sa).Error
	assert.Error(t, err)
}
