package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jrapoport/gothic/storage"
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
	storage.TruncateAll(ts.API.db)

	// Create user
	u, err := models.NewUser("test@example.com", "password", ts.Config.JWT.Aud, nil)
	require.NoError(ts.T(), err, "Error creating test user model")
	require.NoError(ts.T(), ts.API.db.Create(u).Error, "Error saving new test user")
}

func (ts *UserTestSuite) TestUser_UpdatePassword() {
	u, err := models.FindUserByEmailAndAudience(ts.API.db, "test@example.com", ts.Config.JWT.Aud)
	require.NoError(ts.T(), err)

	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"password": "newpass",
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
	require.Equal(ts.T(), w.Code, http.StatusOK)

	u, err = models.FindUserByEmailAndAudience(ts.API.db, "test@example.com", ts.Config.JWT.Aud)
	require.NoError(ts.T(), err)

	assert.True(ts.T(), u.Authenticate("newpass"))
}
