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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const signupEmail = "signup_test@example.com"
const signupPassGood = "test~password"
const signupPassBad = "test"

type SignupTestSuite struct {
	suite.Suite
	api    *API
	config *conf.Configuration
}

func TestSignup(t *testing.T) {
	ts := &SignupTestSuite{}
	suite.Run(t, ts)
}

func (ts *SignupTestSuite) SetupTest() {
	api, config, err := setupAPIForTestForInstance(ts.T())
	ts.api = api
	ts.config = config
	err = ts.api.db.DropDatabase()
	assert.NoError(ts.T(), err)
	ts.config.Webhook = conf.WebhookConfig{}
}

func (ts *SignupTestSuite) parseResponseToken(w *httptest.ResponseRecorder) *GothicClaims {
	res := AccessTokenResponse{}
	require.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&res))
	token := res.Token
	claims := &GothicClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ts.config.JWT.Secret), nil
	},
		jwt.WithAudience(ts.config.JWT.Aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	require.NoError(ts.T(), err, "Error parsing token")
	return claims
}

func (ts *SignupTestSuite) testSignupRequest(params map[string]interface{}) *httptest.ResponseRecorder {
	var buffer bytes.Buffer
	if params != nil {
		require.NoError(ts.T(), json.NewEncoder(&buffer).Encode(params))
	}
	req := httptest.NewRequest(http.MethodPost, "/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.api.handler.ServeHTTP(w, req)
	return w
}

func (ts *SignupTestSuite) TestSignup_Disabled() {
	ts.config.DisableSignup = true
	params := map[string]interface{}{
		"email":    signupEmail,
		"password": signupPassGood,
		"data": map[string]interface{}{
			"a": 1,
		},
	}
	w := ts.testSignupRequest(params)
	assert.Equal(ts.T(), http.StatusForbidden, w.Code)
}

func (ts *SignupTestSuite) TestSignup_Recaptcha() {
	ts.config.Recaptcha.Key = recaptchaDebugKey
	recap := map[string]interface{}{
		"recaptcha": recaptchaDebugToken,
	}
	params := map[string]interface{}{
		"email":    signupEmail,
		"password": signupPassGood,
		"data":     recap,
	}
	w := ts.testSignupRequest(params)
	assert.Equal(ts.T(), http.StatusOK, w.Code)
	recap["recaptcha"] = "nope"
	w = ts.testSignupRequest(params)
	assert.NotEqual(ts.T(), http.StatusOK, w.Code)
}

func (ts *SignupTestSuite) TestSignup_SignupCode() {
	var w *httptest.ResponseRecorder
	ts.config.Signup.Code = true
	params := map[string]interface{}{
		"email":    signupEmail,
		"password": signupPassGood,
	}
	noCode := map[string]interface{}{}
	params["data"] = noCode
	w = ts.testSignupRequest(params)
	assert.NotEqual(ts.T(), http.StatusOK, w.Code)
	badCode := map[string]interface{}{
		"code": "bad",
	}
	params["data"] = badCode
	w = ts.testSignupRequest(params)
	assert.NotEqual(ts.T(), http.StatusOK, w.Code)
	su, err := ts.api.NewSignupCode(models.PINFormat, models.SingleUse)
	assert.NoError(ts.T(), err)
	goodCode := map[string]interface{}{
		"code": su.Code,
	}
	params["data"] = goodCode
	w = ts.testSignupRequest(params)
	assert.Equal(ts.T(), http.StatusOK, w.Code)
	params["email"] = "foo" + signupEmail
	w = ts.testSignupRequest(params)
	assert.NotEqual(ts.T(), http.StatusOK, w.Code)
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
		},
	}))

	// Setup request
	req := httptest.NewRequest(http.MethodPost, "/signup", &buffer)
	req.Header.Set("Content-Type", "application/json")

	// Setup response recorder
	w := httptest.NewRecorder()
	ts.api.handler.ServeHTTP(w, req)
	require.Equal(ts.T(), http.StatusOK, w.Code)
	claims := ts.parseResponseToken(w)
	assert.Equal(ts.T(), signupEmail, claims.Email)
	assert.Equal(ts.T(), 1.0, claims.UserMetaData["a"])
	assert.Equal(ts.T(), "email", claims.AppMetaData["provider"])
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

	ts.api.handler.ServeHTTP(w, req)

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

	ts.api.config.Validation.PasswordRegex = "^[a-z]{2,}$"
	ts.api.handler.ServeHTTP(w, req)
	ts.api.config.Validation.PasswordRegex = ""

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
			return []byte(ts.config.Webhook.Secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
		assert.True(token.Valid)
		assert.Equal(ts.config.JWT.Subject, claims.Subject)
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
		assert.Equal("user", u["role"])
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

	ts.config.Webhook = conf.WebhookConfig{
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
	ts.api.handler.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal(1, callCount)
}

func (ts *SignupTestSuite) TestFailingWebhook() {
	ts.config.Webhook = conf.WebhookConfig{
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

	ts.api.handler.ServeHTTP(w, req)

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

	ts.api.handler.ServeHTTP(y, req)
	u, err := models.FindUserByEmail(ts.api.db, signupEmail)
	if err == nil {
		require.NoError(ts.T(), u.Confirm(ts.api.db))
	}

	encode()
	ts.api.handler.ServeHTTP(w, req)

	data := make(map[string]interface{})
	require.NoError(ts.T(), json.NewDecoder(w.Body).Decode(&data))

	assert.Equal(ts.T(), http.StatusBadRequest, w.Code)
	assert.Equal(ts.T(), float64(http.StatusBadRequest), data["code"])
}

func (ts *SignupTestSuite) TestConfirmSignup() {
	user, err := models.NewUser(signupEmail, "testing", nil)
	user.ConfirmationToken = "asdf3"
	require.NoError(ts.T(), err)
	require.NoError(ts.T(), ts.api.db.Create(user).Error)

	// Find test user
	u, err := models.FindUserByEmail(ts.api.db, signupEmail)
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

	ts.api.handler.ServeHTTP(w, req)

	assert.Equal(ts.T(), http.StatusOK, w.Code, w.Body.String())
}
