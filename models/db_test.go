package models_test

import (
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/storage/test"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"

	"github.com/jrapoport/gothic/models"
	"github.com/stretchr/testify/assert"
)

const modelsTestConfig = "../env/test.env"

func TestTableNameNamespacing(t *testing.T) {
	cases := []struct {
		expected string
		value    interface{}
	}{
		{expected: "test_audit_log_entries", value: []*models.AuditLogEntry{}},
		{expected: "test_refresh_tokens", value: []*models.RefreshToken{}},
		{expected: "test_users", value: []*models.User{}},
	}

	globalConfig, err := conf.LoadConfiguration(modelsTestConfig)
	require.NoError(t, err)

	conn, err := test.SetupDBConnection(t, globalConfig)
	require.NoError(t, err)

	for _, tc := range cases {
		stmt := &gorm.Statement{DB: conn.DB}
		err := stmt.Parse(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, stmt.Schema.Table)
	}
}
