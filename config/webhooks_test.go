package config

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/core/events"
	"github.com/stretchr/testify/assert"
)

const (
	webhookURL    = "http://webhook.example.com/hook"
	webhookSecret = "i-am-a-webhook-secret"
	webhookTry    = 99
	testTimeout   = 100 * time.Minute
)

var hookEvents = []events.Event{events.Login, events.Confirmed}

func TestWebhook(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		w := c.Webhook
		assert.Equal(t, webhookURL+test.mark, w.URL)
		assert.Equal(t, webhookSecret+test.mark, w.Secret)
		for _, event := range w.Events {
			has := w.HasEvent(event)
			assert.True(t, has)
		}
		assert.Equal(t, webhookTry, w.MaxRetries)
		assert.Equal(t, testTimeout, w.Timeout)
		has := w.HasEvent("bad")
		assert.False(t, has)
	})
}

// tests the ENV vars are correctly taking precedence
func TestWebhook_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			w := c.Webhook
			assert.Equal(t, webhookURL, w.URL)
			assert.Equal(t, webhookSecret, w.Secret)
			assert.Equal(t, hookEvents, w.Events)
			assert.Equal(t, webhookTry, w.MaxRetries)
			assert.Equal(t, testTimeout, w.Timeout)
			for _, event := range hookEvents {
				has := w.HasEvent(event)
				assert.True(t, has)
			}
		})
	}
}

// test the *un-normalized* defaults with load
func TestWebhook_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	def := networkDefaults
	n := c.Network
	assert.Equal(t, def, n)
}

func TestWebhook_Normalization(t *testing.T) {
	w := Webhooks{}
	j := jwtDefaults
	j.Secret = webhookSecret
	err := w.normalize(serviceDefaults, j)
	assert.NoError(t, err)
	assert.False(t, w.Enabled())
	assert.Equal(t, "", w.Secret)
	w.URL = webhookURL
	err = w.normalize(Service{
		Name:    service,
		SiteURL: siteURL,
	}, JWT{
		Secret: webhookSecret,
	})
	assert.NoError(t, err)
	assert.True(t, w.Enabled())
	assert.Equal(t, webhookSecret, w.Secret)
	w = Webhooks{URL: "\n"}
	err = w.normalize(serviceDefaults, j)
	assert.Error(t, err)
	w = Webhooks{
		URL: webhookURL,
		Events: []events.Event{
			events.Login,
			events.All,
			events.Signup,
		},
	}
	err = w.normalize(serviceDefaults, j)
	assert.NoError(t, err)
	assert.Equal(t, []events.Event{
		events.All,
	}, w.Events)
}

func TestFormatWebhookURL(t *testing.T) {
	u := FormatWebhookURL(webhookURL, events.Login)
	assert.Equal(t, webhookURL, u)
	u = FormatWebhookURL(webhookURL+"/"+WebhookURLEvent, events.Login)
	assert.Equal(t, webhookURL+"/"+string(events.Login), u)
}
