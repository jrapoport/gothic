package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const confirmEmail = "confirm_test@example.com"

type ConfirmTestSuite struct {
	suite.Suite
	a *API
	c *conf.Configuration
}

func TestConfirm(t *testing.T) {
	api, config, err := setupAPIForTestForInstance(t)
	require.NoError(t, err)

	ts := &ConfirmTestSuite{
		a: api,
		c: config,
	}

	suite.Run(t, ts)
}

func (ts *ConfirmTestSuite) SetupTest() {
	err := ts.a.db.DropDatabase()
	assert.NoError(ts.T(), err)
	// Create user
	u, err := models.NewUser(confirmEmail, "password", nil)
	require.NoError(ts.T(), err, "Error creating test user model")
	require.NoError(ts.T(), ts.a.db.Create(u).Error, "Error saving new test user")
}

func (ts *ConfirmTestSuite) TestConfirm_PasswordRecovery() {
	u, err := models.FindUserByEmail(ts.a.db, confirmEmail)
	require.NoError(ts.T(), err)
	u.RecoverySentAt = &time.Time{}
	require.NoError(ts.T(), ts.a.db.Save(u).Error)

	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email": confirmEmail,
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "http://localhost/recover", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()
	ts.a.handler.ServeHTTP(w, req)
	assert.Equal(ts.T(), http.StatusOK, w.Code)

	u, err = models.FindUserByEmail(ts.a.db, confirmEmail)
	require.NoError(ts.T(), err)

	assert.WithinDuration(ts.T(), time.Now(), *u.RecoverySentAt, 1*time.Second)
	assert.False(ts.T(), u.IsConfirmed())

	// Send Confirm request
	var vbuffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&vbuffer).Encode(map[string]interface{}{
		"type":  "recovery",
		"token": u.RecoveryToken,
	}))

	req = httptest.NewRequest(http.MethodPost, "http://localhost/verify", &vbuffer)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	ts.a.handler.ServeHTTP(w, req)
	assert.Equal(ts.T(), http.StatusOK, w.Code)

	u, err = models.FindUserByEmail(ts.a.db, confirmEmail)
	require.NoError(ts.T(), err)
	assert.True(ts.T(), u.IsConfirmed())
}
