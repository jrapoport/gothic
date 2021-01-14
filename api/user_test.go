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
	"github.com/jrapoport/gothic/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const userEmail = "user_test@example.com"

type UserTestSuite struct {
	suite.Suite
	API    *API
	Config *conf.Configuration
}

func TestUser(t *testing.T) {
	api, config, err := setupAPIForTestForInstance(t)
	require.NoError(t, err)

	ts := &UserTestSuite{
		API:    api,
		Config: config,
	}

	suite.Run(t, ts)
}

func (ts *UserTestSuite) SetupTest() {
	// Create user
	u, err := models.NewUser(userEmail, "password", ts.Config.JWT.Aud, nil)
	require.NoError(ts.T(), err, "Error creating test user model")
	t := time.Now()
	u.ConfirmedAt = &t
	require.NoError(ts.T(), ts.API.db.Create(u).Error, "Error saving new test user")
}

func (ts *UserTestSuite) TearDownTest() {
	storage.TruncateAll(ts.API.db)
}

func (ts *UserTestSuite) TestUser_UpdatePassword() {
	const password = "new!password"
	u, err := models.FindUserByEmailAndAudience(ts.API.db, userEmail, ts.Config.JWT.Aud)
	require.NoError(ts.T(), err)

	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"password": password,
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPut, "http://localhost/user", &buffer)
	req.Header.Set("Content-Type", "application/json")

	token, err := generateAccessToken(u,
		time.Second*time.Duration(ts.Config.JWT.Exp),
		ts.Config.JWT.Secret,
		ts.Config.JWT.SigningMethod())
	require.NoError(ts.T(), err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Setup response recorder
	w := httptest.NewRecorder()
	ts.API.handler.ServeHTTP(w, req)
	require.Equal(ts.T(), http.StatusOK, w.Code)

	u, err = models.FindUserByEmailAndAudience(ts.API.db, userEmail, ts.Config.JWT.Aud)
	require.NoError(ts.T(), err)

	assert.True(ts.T(), u.Authenticate(password))
}
