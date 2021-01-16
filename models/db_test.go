package models_test

import (
	"gorm.io/gorm"
	"testing"

	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"github.com/jrapoport/gothic/storage/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const modelsTestConfig = "../env/test.env"

func TestTableNameNamespacing(t *testing.T) {

	config, err := conf.LoadConfiguration(modelsTestConfig)
	require.NoError(t, err)
	ns := storage.Namespace(config)
	cases := []struct {
		expected string
		value    interface{}
	}{
		{expected: models.AuditLogEntry{}.TableName(), value: []*models.AuditLogEntry{}},
		{expected: ns + "refresh_tokens", value: []*models.RefreshToken{}},
		{expected: ns + "signup_codes", value: []*models.SignupCode{}},
		{expected: ns + "users", value: []*models.User{}},
	}

	conn, err := test.SetupDBConnection(t, config)
	require.NoError(t, err)

	for _, tc := range cases {
		stmt := &gorm.Statement{DB: conn.DB}
		err = stmt.Parse(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, stmt.Schema.Table)
	}
}
