package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

const userEmail = "user_test@example.com"

type UserTestSuite struct {
	suite.Suite
	a *API
	c *conf.Configuration
}

func TestUser(t *testing.T) {
	api, config, err := setupAPIForTestForInstance(t)
	require.NoError(t, err)

	ts := &UserTestSuite{
		a: api,
		c: config,
	}

	suite.Run(t, ts)
}

func (ts *UserTestSuite) SetupTest() {
	// Create user
	u, err := models.NewUser(userEmail, "password", nil)
	require.NoError(ts.T(), err, "Error creating test user model")
	t := time.Now()
	u.ConfirmedAt = &t
	require.NoError(ts.T(), ts.a.db.Create(u).Error, "Error saving new test user")
}

func (ts *UserTestSuite) TearDownTest() {
	err := ts.a.db.DropDatabase()
	assert.NoError(ts.T(), err)
}

func (ts *UserTestSuite) TestUser_UpdatePassword() {
	const password = "new!password"
	u, err := models.FindUserByEmail(ts.a.db, userEmail)
	require.NoError(ts.T(), err)

	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"password": password,
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPut, "http://localhost/user", &buffer)
	req.Header.Set("Content-Type", "application/json")

	token, err := generateAccessToken(u, ts.c.JWT)
	require.NoError(ts.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Setup response recorder
	w := httptest.NewRecorder()
	ts.a.handler.ServeHTTP(w, req)
	require.Equal(ts.T(), http.StatusOK, w.Code)

	u, err = models.FindUserByEmail(ts.a.db, userEmail)
	require.NoError(ts.T(), err)

	assert.True(ts.T(), u.Authenticate(password))
}
