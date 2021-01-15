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
	"github.com/jrapoport/gothic/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const confirmEmail = "confirm_test@example.com"

type ConfirmTestSuite struct {
	suite.Suite
	API    *API
	Config *conf.Configuration
}

func TestConfirm(t *testing.T) {
	api, config, err := setupAPIForTestForInstance(t)
	require.NoError(t, err)

	ts := &ConfirmTestSuite{
		API:    api,
		Config: config,
	}

	suite.Run(t, ts)
}

func (ts *ConfirmTestSuite) SetupTest() {
	storage.TruncateAll(ts.API.db)

	// Create user
	u, err := models.NewUser(confirmEmail, "password", nil)
	require.NoError(ts.T(), err, "Error creating test user model")
	require.NoError(ts.T(), ts.API.db.Create(u).Error, "Error saving new test user")
}

func (ts *ConfirmTestSuite) TestConfirm_PasswordRecovery() {
	u, err := models.FindUserByEmail(ts.API.db, confirmEmail)
	require.NoError(ts.T(), err)
	u.RecoverySentAt = &time.Time{}
	require.NoError(ts.T(), ts.API.db.Save(u).Error)

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
	ts.API.handler.ServeHTTP(w, req)
	assert.Equal(ts.T(), http.StatusOK, w.Code)

	u, err = models.FindUserByEmail(ts.API.db, confirmEmail)
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
	ts.API.handler.ServeHTTP(w, req)
	assert.Equal(ts.T(), http.StatusOK, w.Code)

	u, err = models.FindUserByEmail(ts.API.db, confirmEmail)
	require.NoError(ts.T(), err)
	assert.True(ts.T(), u.IsConfirmed())
}
