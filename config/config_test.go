package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDir    = "testdata/"
	testDotEnv = testDir + "config_test.dot.env"
	testEnv    = testDir + "config_test.env"
	testYAML   = testDir + "config_test.yaml"
	testJSON   = testDir + "config_test.json"
)

func clearEnv() {
	os.Clearenv()
}

func setEnv(t *testing.T, k, v string) {
	err := os.Setenv(k, v)
	require.NoError(t, err)
}

func loadDotEnv(t *testing.T) {
	clearEnv()
	err := godotenv.Load(testDotEnv)
	require.NoError(t, err)
}

func loadCleanEnv(*testing.T) {
	clearEnv()
}

type testCase struct {
	name    string
	file    string
	mark    string
	testEnv func(t *testing.T)
}

var tests = []testCase{
	{"ENV", "", "", loadDotEnv},
	{"dotenv", testEnv, filepath.Ext(testEnv), loadCleanEnv},
	{"yaml", testYAML, filepath.Ext(testYAML), loadCleanEnv},
	{"json", testJSON, filepath.Ext(testJSON), loadCleanEnv},
}

func runTests(t *testing.T, runTest func(t *testing.T, test testCase, c *Config)) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.testEnv(t)
			test.mark = filepath.Ext(test.file)
			c := loadTestConfig(t, test)
			runTest(t, test, c)
		})
	}
}

func loadTestConfig(t *testing.T, test testCase) *Config {
	test.testEnv(t)
	c, err := loadNormalized(test.file)
	require.NoError(t, err)
	require.NotNil(t, c)
	return c
}

func TestLoadConfig(t *testing.T) {
	reqTests := []struct {
		key   string
		value string
		Err   assert.ErrorAssertionFunc
	}{
		{ENVPrefix + "_SITE_URL", siteURL, assert.Error},
		{ENVPrefix + "_DB_DSN", dsn, assert.Error},
		{ENVPrefix + "_JWT_SECRET", jwtSecret, assert.Error},
	}
	reqTests[len(reqTests)-1].Err = assert.NoError
	clearEnv()
	for _, test := range reqTests {
		setEnv(t, test.key, test.value)
		_, err := LoadConfig("")
		test.Err(t, err)
	}
	_, err := LoadConfig(testEnv)
	assert.Error(t, err)
	setEnv(t, ENVPrefix+"LOG_LEVEL", "bad")
	_, err = LoadConfig("")
	assert.Error(t, err)
	_, err = LoadConfig("bad path")
	assert.Error(t, err)
	_, err = LoadConfig("\n.json")
	assert.Error(t, err)
}

func TestLog(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		l := c.Log()
		c.Log().Info("test")
		c.ReplaceLog(logrus.New())
		c.Log().Info("test")
		c.ReplaceLog(l)
	})
}

// test the *un-normalized* defaults with load
func TestConfig_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	def := configDefaults
	assert.Equal(t, def, *c)
}

func TestConfig_Normalization(t *testing.T) {
	c := Config{}
	err := c.normalize()
	assert.Error(t, err)
	c = Config{
		Service: Service{
			SiteURL: siteURL,
		},
		DB: Database{
			DSN: dsn,
		},
		Security: Security{
			JWT: JWT{
				Secret: jwtSecret,
			},
		},
	}
	err = c.normalize()
	assert.Error(t, err)
	assert.Equal(t, BuildVersion(), c.Version())
}

func TestWriteConfig(t *testing.T) {
	dir := t.TempDir()
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		_, name := filepath.Split(test.file)
		file := filepath.Join(dir, "cfg"+name)
		err := c.Write(file)
		assert.NoError(t, err)
	})
}
