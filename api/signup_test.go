package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const signupEmail = "signup_test@example.com"
const signupPassGood = "test~password"
const signupPassBad = "test"

type SignupTestSuite struct {
	suite.Suite
	API    *API
	Config *conf.Configuration
}

func TestSignup(t *testing.T) {
	api, config, err := setupAPIForTestForInstance(t)
	require.NoError(t, err)

	ts := &SignupTestSuite{
		API:    api,
		Config: config,
	}
	// defer api.db.Close()

	suite.Run(t, ts)
}

func (ts *SignupTestSuite) SetupTest() {
	storage.TruncateAll(ts.API.db)
	ts.Config.Webhook = conf.WebhookConfig{}
}

// TestSignup tests API /signup route
func (ts *SignupTestSuite) TestSignup() {
	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email":    signupEmail,
		"password": signupPassGood,
		"data": map[string]interface{}{
			"a": 1,
			"recaptcha": "hello",
		},
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.API.handler.ServeHTTP(w, req)

	require.Equal(ts.T(), http.StatusOK, w.Code)

	data := models.User{}
	require.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&data))
	assert.Equal(ts.T(), signupEmail, data.Email)
	assert.Equal(ts.T(), ts.Config.JWT.Aud, data.Aud)
	assert.Equal(ts.T(), 1.0, data.UserMetaData["a"])
	assert.Equal(ts.T(), "email", data.AppMetaData["provider"])
}

// TestSignup tests API /signup route
func (ts *SignupTestSuite) TestSignup_BadPassword() {
	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email":    signupEmail,
		"password": signupPassBad,
		"data": map[string]interface{}{
			"a": 1,
		},
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.API.handler.ServeHTTP(w, req)

	require.NotEqual(ts.T(), http.StatusOK, w.Code)
}

// TestSignup tests API /signup route
func (ts *SignupTestSuite) TestSignup_PasswordRegex() {
	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email":    signupEmail,
		"password": signupPassBad,
		"data": map[string]interface{}{
			"a": 1,
		},
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.API.config.PasswordRegex = "^[a-z]{2,}$"
	ts.API.handler.ServeHTTP(w, req)
	ts.API.config.PasswordRegex = ""

	require.Equal(ts.T(), http.StatusOK, w.Code)
}

func (ts *SignupTestSuite) TestWebhookTriggered() {
	var callCount int
	require := ts.Require()
	assert := ts.Assert()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		// verify the signature
		signature := r.Header.Get("x-webhook-signature")
		claims := new(jwt.StandardClaims)
		token, err := jwt.ParseWithClaims(signature, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(ts.Config.Webhook.Secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
		assert.True(token.Valid)
		assert.Equal(ts.Config.JWT.Subject, claims.Subject) // not configured for multitenancy
		assert.Equal("gothic", claims.Issuer)
		assert.WithinDuration(time.Now(), claims.IssuedAt.Time, 5*time.Second)

		// verify the contents
		defer squash(r.Body.Close)
		raw, err := ioutil.ReadAll(r.Body)
		require.NoError(err)
		data := map[string]interface{}{}
		require.NoError(json.Unmarshal(raw, &data))

		assert.Equal(2, len(data))
		assert.Equal("validate", data["event"])

		u, ok := data["user"].(map[string]interface{})
		require.True(ok)
		assert.Len(u, 10)
		// assert.Equal(t, user.ID, u["id"]) TODO
		assert.Equal("api.gothic.com", u["aud"])
		assert.Equal("", u["role"])
		assert.Equal(signupEmail, u["email"])

		appmeta, ok := u["app_metadata"].(map[string]interface{})
		require.True(ok)
		assert.Len(appmeta, 1)
		assert.EqualValues("email", appmeta["provider"])

		usermeta, ok := u["user_metadata"].(map[string]interface{})
		require.True(ok)
		assert.Len(usermeta, 1)
		assert.EqualValues(1, usermeta["a"])
	}))

	// Allowing connection to localhost for the tests only
	localhost := removeLocalhostFromPrivateIPBlock()
	defer unshiftPrivateIPBlock(localhost)

	ts.Config.Webhook = conf.WebhookConfig{
		URL:        svr.URL,
		Retries:    1,
		TimeoutSec: 1,
		Secret:     "top-secret",
		Events:     []string{"validate"},
	}
	var buffer bytes.Buffer
	require.NoError(json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email":    signupEmail,
		"password": signupPassGood,
		"data": map[string]interface{}{
			"a": 1,
		},
	}))
	req := httptest.NewRequest(http.MethodPost, "http://localhost/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ts.API.handler.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal(1, callCount)
}

func (ts *SignupTestSuite) TestFailingWebhook() {
	ts.Config.Webhook = conf.WebhookConfig{
		URL:        "http://notaplace.localhost",
		Retries:    1,
		TimeoutSec: 1,
		Events:     []string{"validate", "signup"},
	}
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"email":    signupEmail,
		"password": signupPassGood,
		"data": map[string]interface{}{
			"a": 1,
		},
	}))
	req := httptest.NewRequest(http.MethodPost, "http://localhost/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.API.handler.ServeHTTP(w, req)

	require.Equal(ts.T(), http.StatusBadGateway, w.Code)
}

// TestSignupTwice checks to make sure the same email cannot be registered twice
func (ts *SignupTestSuite) TestSignupTwice() {
	// Request body
	var buffer bytes.Buffer

	encode := func() {
		require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
			"email":    signupEmail,
			"password": signupPassGood,
			"data": map[string]interface{}{
				"a": 1,
			},
		}))
	}

	encode()

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "http://localhost/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()
	y := httptest.NewRecorder()

	ts.API.handler.ServeHTTP(y, req)
	u, err := models.FindUserByEmailAndAudience(ts.API.db, signupEmail, ts.Config.JWT.Aud)
	if err == nil {
		require.NoError(ts.T(), u.Confirm(ts.API.db))
	}

	encode()
	ts.API.handler.ServeHTTP(w, req)

	data := make(map[string]interface{})
	require.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&data))

	assert.Equal(ts.T(), http.StatusBadRequest, w.Code)
	assert.Equal(ts.T(), float64(http.StatusBadRequest), data["code"])
}

func (ts *SignupTestSuite) TestConfirmSignup() {
	user, err := models.NewUser(signupEmail, "testing", ts.Config.JWT.Aud, nil)
	user.ConfirmationToken = "asdf3"
	require.NoError(ts.T(), err)
	require.NoError(ts.T(), ts.API.db.Create(user).Error)

	// Find test user
	u, err := models.FindUserByEmailAndAudience(ts.API.db, signupEmail, ts.Config.JWT.Aud)
	require.NoError(ts.T(), err)

	// Request body
	var buffer bytes.Buffer
	require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(map[string]interface{}{
		"type":  "signup",
		"token": u.ConfirmationToken,
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "http://localhost/verify", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()

	ts.API.handler.ServeHTTP(w, req)

	assert.Equal(ts.T(), http.StatusOK, w.Code, w.Body.String())
}
