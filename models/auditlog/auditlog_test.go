package auditlog

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/store/migration"
	"github.com/jrapoport/gothic/test/tconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func auditConn(t *testing.T) *store.Connection {
	conn, _ := tconn.TempConn(t)
	mg := migration.NewMigration("1", AuditLog{})
	err := conn.RunMigration(mg)
	require.NoError(t, err)
	return conn
}

func TestNewAuditLog(t *testing.T) {
	conn := auditConn(t)
	l := NewAuditLog(Token, Banned, uuid.New(), types.Map{key.Token: "token"})
	err := conn.Create(l).Error
	assert.NoError(t, err)
}
