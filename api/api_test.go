package api

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/storage"
	"github.com/jrapoport/gothic/storage/test"
	"github.com/stretchr/testify/require"
)

const (
	apiTestVersion = "1"
	apiTestConfig  = "../env/test.env"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// setupAPIForTest creates a new API to run tests with.
// Using this function allows us to keep track of the database connection
// and cleaning up data between tests.
func setupAPIForTest() (*API, *conf.Configuration, error) {
	return setupAPIForTestWithCallback(nil)
}

func setupAPIForTestForInstance() (*API, *conf.Configuration, error) {
	// BUG: is this right? seems ok to ditch this.
	/*
		cb := func(gc *conf.GlobalConfiguration, c *conf.Configuration, conn *storage.Connection) error {
			err := conn.Create(&models.Instance{
				BaseConfig: c,
			}).Error
			return err
		}

		api, conf, err := setupAPIForTestWithCallback(cb)
	*/
	api, c, err := setupAPIForTestWithCallback(nil)
	if err != nil {
		return nil, nil, err
	}
	return api, c, nil
}

func setupAPIForTestWithCallback(cb func(*conf.GlobalConfiguration, *conf.Configuration, *storage.Connection) error) (*API, *conf.Configuration, error) {
	globalConfig, err := conf.LoadGlobal(apiTestConfig)
	if err != nil {
		return nil, nil, err
	}

	conn, err := test.SetupDBConnection(globalConfig)
	if err != nil {
		return nil, nil, err
	}

	config, err := conf.LoadConfig(apiTestConfig)
	if err != nil {
		return nil, nil, err
	}

	if cb != nil {
		err = cb(globalConfig, config, conn)
		if err != nil {
			return nil, nil, err
		}
	}

	ctx, err := WithConfig(context.Background(), config)
	if err != nil {
		return nil, nil, err
	}

	return NewAPIWithVersion(ctx, globalConfig, conn, apiTestVersion), config, nil
}

func TestEmailEnabledByDefault(t *testing.T) {
	api, _, err := setupAPIForTest()
	require.NoError(t, err)

	require.False(t, api.config.External.Email.Disabled)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
