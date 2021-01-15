package api

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuditTestSuite struct {
	suite.Suite
	API    *API
	Config *conf.Configuration

	token string
}

const auditAdminEmail = "admin@audit.com"
const auditUserEmail = "user@audit.com"

func TestAudit(t *testing.T) {
	api, config, err := setupAPIForTestForInstance(t)
	require.NoError(t, err)

	ts := &AuditTestSuite{
		API:    api,
		Config: config,
	}

	suite.Run(t, ts)
}

func (ts *AuditTestSuite) SetupTest() {
	storage.TruncateAll(ts.API.db)
	ts.token = ts.makeSuperAdmin(auditAdminEmail)
}

func (ts *AuditTestSuite) makeSuperAdmin(email string) string {
	u, err := models.NewUser(email, "test", map[string]interface{}{"full_name": "Test Username"})
	require.NoError(ts.T(), err, "Error making new user")

	u.IsSuperAdmin = true
	require.NoError(ts.T(), ts.API.db.Create(u).Error, "Error creating user")

	token, err := generateAccessToken(u, ts.Config.JWT)
	require.NoError(ts.T(), err, "Error generating access token")

	_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(ts.Config.JWT.Secret), nil
	},
		jwt.WithAudience(ts.Config.JWT.Aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	require.NoError(ts.T(), err, "Error parsing token")

	return token
}

func (ts *AuditTestSuite) TestAuditGet() {
	ts.prepareDeleteEvent()
	// CHECK FOR AUDIT LOG
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/audit", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ts.token))

	ts.API.handler.ServeHTTP(w, req)
	require.Equal(ts.T(), http.StatusOK, w.Code)

	assert.Equal(ts.T(), "</admin/audit?page=1>; rel=\"last\"", w.HeaderMap.Get("Link"))
	assert.Equal(ts.T(), "1", w.HeaderMap.Get("X-Total-Count"))

	logs := []models.AuditLogEntry{}
	require.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&logs))

	require.Len(ts.T(), logs, 1)
	require.Contains(ts.T(), logs[0].Payload, "actor_email")
	assert.Equal(ts.T(), auditAdminEmail, logs[0].Payload["actor_email"])
	traits, ok := logs[0].Payload["traits"].(map[string]interface{})
	require.True(ts.T(), ok)
	require.Contains(ts.T(), traits, "user_email")
	assert.Equal(ts.T(), auditUserEmail, traits["user_email"])
}

func (ts *AuditTestSuite) TestAuditFilters() {
	ts.prepareDeleteEvent()

	queries := []string{
		"/admin/audit?query=action:user_deleted",
		"/admin/audit?query=type:team",
		"/admin/audit?query=author:user",
		"/admin/audit?query=author:@audit.com",
	}

	for _, q := range queries {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, q, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ts.token))

		ts.API.handler.ServeHTTP(w, req)
		require.Equal(ts.T(), http.StatusOK, w.Code)
		if w.Code != http.StatusOK {
			ts.T().Log(w.Body.String())
			ts.T().FailNow()
		}

		logs := []models.AuditLogEntry{}
		require.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&logs))

		require.Len(ts.T(), logs, 1)
		require.Contains(ts.T(), logs[0].Payload, "actor_email")
		assert.Equal(ts.T(), auditAdminEmail, logs[0].Payload["actor_email"])
		traits, ok := logs[0].Payload["traits"].(map[string]interface{})
		require.True(ts.T(), ok)
		require.Contains(ts.T(), traits, "user_email")
		assert.Equal(ts.T(), auditUserEmail, traits["user_email"])
	}
}

func (ts *AuditTestSuite) prepareDeleteEvent() {
	// DELETE USER
	u, err := models.NewUser(auditUserEmail, "test", nil)
	require.NoError(ts.T(), err, "Error making new user")
	require.NoError(ts.T(), ts.API.db.Create(u).Error, "Error creating user")

	// Setup request
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/users/%s", u.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ts.token))

	ts.API.handler.ServeHTTP(w, req)
	require.Equal(ts.T(), http.StatusOK, w.Code)
}
