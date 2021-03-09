package config

import (
	"testing"

	"github.com/jrapoport/gothic/store/drivers"
	"github.com/stretchr/testify/assert"
)

const (
	namespace       = "foo"
	driver          = drivers.MySQL
	dsn             = "root@tcp(0.0.0.0:3306)/test?parseTime=true&foo=a"
	maxRetries  int = 99
	autoMigrate     = true
)

func TestDatabase(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		db := c.DB
		assert.Equal(t, namespace+test.mark, db.Namespace)
		assert.EqualValues(t, driver, db.Driver)
		assert.Equal(t, dsn+test.mark, db.DSN)
		assert.Equal(t, maxRetries, db.MaxRetries)
		assert.Equal(t, autoMigrate, db.AutoMigrate)
	})
}

// tests the ENV vars are correctly taking precedence
func TestDatabase_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			db := c.DB
			assert.Equal(t, namespace, db.Namespace)
			assert.EqualValues(t, driver, db.Driver)
			assert.Equal(t, dsn, db.DSN)
			assert.Equal(t, maxRetries, db.MaxRetries)
			assert.Equal(t, autoMigrate, db.AutoMigrate)
		})
	}
}

// test the *un-normalized* defaults with load
func TestDatabase_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	def := databaseDefaults
	db := c.DB
	assert.Equal(t, def, db)
}
