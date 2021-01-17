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

type RecoverTestSuite struct {
	suite.Suite
	a *API
	c *conf.Configuration
}

func TestRecover(t *testing.T) {
	api, config, err := setupAPIForTestForInstance(t)
	require.NoError(t, err)

	ts := &RecoverTestSuite{
		a: api,
		c: config,
	}

	suite.Run(t, ts)
}

func (ts *RecoverTestSuite) SetupTest() {
	err := ts.a.db.DropDatabase()
	assert.NoError(ts.T(), err)
	// Create user
	u, err := models.NewUser("test@example.com", "password", nil)
	require.NoError(ts.T(), err, "Error creating test user model")
	require.NoError(ts.T(), ts.a.db.Create(u).Error, "Error saving new test user")
}

func (ts *RecoverTestSuite) TestRecover_FirstRecovery() {
	u, err := models.FindUserByEmail(ts.a.db, "test@example.com")
	require.NoError(ts.T(), err)
	u.RecoverySentAt = &time.Time{}
	require.NoError(ts.T(), ts.a.db.Save(u).Error)

	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email": "test@example.com",
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "http://localhost/recover", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()
	ts.a.handler.ServeHTTP(w, req)
	assert.Equal(ts.T(), http.StatusOK, w.Code)

	u, err = models.FindUserByEmail(ts.a.db, "test@example.com")
	require.NoError(ts.T(), err)

	assert.WithinDuration(ts.T(), time.Now(), *u.RecoverySentAt, 1*time.Second)
}

func (ts *RecoverTestSuite) TestRecover_NoEmailSent() {
	recoveryTime := time.Now().UTC().Add(-5 * time.Minute)
	u, err := models.FindUserByEmail(ts.a.db, "test@example.com")
	require.NoError(ts.T(), err)
	u.RecoverySentAt = &recoveryTime
	require.NoError(ts.T(), ts.a.db.Save(u).Error)

	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email": "test@example.com",
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "http://localhost/recover", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()
	ts.a.handler.ServeHTTP(w, req)
	assert.Equal(ts.T(), http.StatusOK, w.Code)

	u, err = models.FindUserByEmail(ts.a.db, "test@example.com")
	require.NoError(ts.T(), err)

	// ensure it did not send a new email
	u1 := recoveryTime.Round(time.Second).Unix()
	u2 := u.RecoverySentAt.Round(time.Second).Unix()
	assert.Equal(ts.T(), u1, u2)
}

func (ts *RecoverTestSuite) TestRecover_NewEmailSent() {
	recoveryTime := time.Now().UTC().Add(-20 * time.Minute)
	u, err := models.FindUserByEmail(ts.a.db, "test@example.com")
	require.NoError(ts.T(), err)
	u.RecoverySentAt = &recoveryTime
	require.NoError(ts.T(), ts.a.db.Save(u).Error)

	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email": "test@example.com",
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "http://localhost/recover", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()
	ts.a.handler.ServeHTTP(w, req)
	assert.Equal(ts.T(), http.StatusOK, w.Code)

	u, err = models.FindUserByEmail(ts.a.db, "test@example.com")
	require.NoError(ts.T(), err)

	// ensure it sent a new email
	assert.WithinDuration(ts.T(), time.Now(), *u.RecoverySentAt, 1*time.Second)
}
