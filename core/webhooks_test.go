package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/core/webhooks"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
)

func testHook(t *testing.T) *API {
	const (
		webhookSecret = "test-secret"
		webhookIssuer = "test"
		webhookFormat = "/hook/%s?event=%s"
	)
	var hookURL = fmt.Sprintf(webhookFormat,
		config.WebhookURLEvent,
		config.WebhookURLEvent)
	a := apiWithTempDB(t)
	c := a.config
	c.Signup.AutoConfirm = false
	c.Validation.UsernameRegex = ""
	c.Validation.PasswordRegex = ""
	c.Webhook.Secret = webhookSecret
	c.Webhook.Issuer = webhookIssuer
	c.Webhook.URL = hookURL
	c.Webhook.Timeout = 0
	c.Webhook.MaxRetries = 0
	return a
}

func testWebhookEvent(t *testing.T, evt events.Event, test events.Event) {
	var mu sync.RWMutex
	a := testHook(t)
	c := a.config
	c.Webhook.Events = []events.Event{test}
	srv, rec := testHookSvr(t, &mu, c.Webhook, evt)
	c.Webhook.URL = srv.URL + c.Webhook.URL
	err := a.LoadWebhooks()
	assert.NoError(t, err)
	email := tutils.RandomEmail()
	ctx := testContext(a)
	_, err = a.Signup(ctx, email, "", "", nil)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return rec.Code == http.StatusOK
	}, 200*time.Millisecond, 10*time.Millisecond)
	var msg types.Map
	err = json.Unmarshal(rec.Body.Bytes(), &msg)
	assert.NoError(t, err)
	assert.EqualValues(t, evt, msg[key.Event])
	assert.EqualValues(t, c.Provider(), msg[key.Provider])
	assert.Equal(t, email, msg[key.Email])
}

func testWebhookSignature(t *testing.T, j config.JWT, r *http.Request) {
	sig := r.Header.Get(webhooks.WebhookSignature)
	var claims webhooks.WebhookClaims
	err := jwt.ParseClaims(j, sig, &claims)
	assert.NoError(t, err)
}

func testHookSvr(t *testing.T, mu *sync.RWMutex, c config.Webhooks, evt events.Event) (*httptest.Server, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	// recorder defaults to OK
	rec.Code = http.StatusNotFound
	record := func(r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		err = r.Body.Close()
		assert.NoError(t, err)
		_, err = rec.Write(body)
		assert.NoError(t, err)
		rec.WriteHeader(http.StatusOK)
		assert.NoError(t, err)
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		u := config.FormatWebhookURL(c.URL, evt)
		assert.Equal(t, u, r.URL.String())
		testWebhookSignature(t, c.JWT, r)
		if c.Timeout > 0 {
			<-time.After(2 * c.Timeout)
			return
		}
		if c.MaxRetries > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			c.MaxRetries--
			return
		}
		record(r)
		w.WriteHeader(http.StatusOK)
	}
	srv := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(func() {
		srv.Close()
	})
	return srv, rec
}

func TestAPI_Webhook(t *testing.T) {
	t.Parallel()
	testWebhookEvent(t, events.Signup, events.Signup)
}

func TestAPI_WebhookAll(t *testing.T) {
	t.Parallel()
	testWebhookEvent(t, events.Signup, events.All)
}

func TestAPI_WebhookTimeout(t *testing.T) {
	t.Parallel()
	var mu sync.RWMutex
	a := testHook(t)
	c := a.config
	evt := events.Signup
	c.Webhook.Events = []events.Event{evt}
	c.Webhook.Timeout = 100 * time.Millisecond
	srv, rec := testHookSvr(t, &mu, c.Webhook, evt)
	c.Webhook.URL = srv.URL + c.Webhook.URL
	err := a.LoadWebhooks()
	assert.NoError(t, err)
	email := tutils.RandomEmail()
	ctx := testContext(a)
	_, err = a.Signup(ctx, email, "", "", nil)
	assert.NoError(t, err)
	assert.Never(t, func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return rec.Code == http.StatusOK
	}, c.Webhook.Timeout*2, 10*time.Millisecond)
}

func TestAPI_WebhookRetry(t *testing.T) {
	t.Parallel()
	var mu sync.RWMutex
	a := testHook(t)
	c := a.config
	evt := events.Signup
	c.Webhook.Events = []events.Event{evt}
	c.Webhook.MaxRetries = 2
	srv, rec := testHookSvr(t, &mu, c.Webhook, evt)
	c.Webhook.URL = srv.URL + c.Webhook.URL
	err := a.LoadWebhooks()
	assert.NoError(t, err)
	email := tutils.RandomEmail()
	ctx := testContext(a)
	_, err = a.Signup(ctx, email, "", "", nil)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return rec.Code == http.StatusOK
	}, 2*time.Second, 100*time.Millisecond)
}

func TestAPI_WebhookNotFound(t *testing.T) {
	t.Parallel()
	var mu sync.RWMutex
	a := testHook(t)
	c := a.config
	evt := events.Signup
	c.Webhook.Events = []events.Event{evt}
	c.Webhook.MaxRetries = 2
	_, rec := testHookSvr(t, &mu, c.Webhook, evt)
	c.Webhook.URL = "http://127.0.0.1:1" + c.Webhook.URL
	err := a.LoadWebhooks()
	assert.NoError(t, err)
	email := tutils.RandomEmail()
	ctx := testContext(a)
	_, err = a.Signup(ctx, email, "", "", nil)
	assert.NoError(t, err)
	assert.Never(t, func() bool {
		mu.RLock()
		defer mu.RUnlock()
		return rec.Code == http.StatusOK
	}, 2*time.Second, 10*time.Millisecond)
}

func TestAPI_WebhookDisabled(t *testing.T) {
	t.Parallel()
	a := testHook(t)
	c := a.config
	// no events
	err := a.LoadWebhooks()
	assert.NoError(t, err)
	// bad url
	c.Webhook.URL = "\n"
	err = a.LoadWebhooks()
	assert.Error(t, err)
	// no secret
	c.Webhook.Secret = ""
	err = a.LoadWebhooks()
	assert.Error(t, err)
	// disabled
	c.Webhook.URL = ""
	err = a.LoadWebhooks()
	assert.NoError(t, err)
}
