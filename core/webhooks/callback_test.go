package webhooks

import (
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/store/types"
	"github.com/stretchr/testify/assert"
)

func TestNewCallback(t *testing.T) {
	const (
		testURL      = "http://example.com/" + config.WebhookURLEvent
		testCallback = "http://example.com/signup"
	)
	c := config.Webhooks{}
	_, err := NewCallback(c, events.Unknown, nil)
	assert.Error(t, err)
	_, err = NewCallback(c, events.Signup, nil)
	assert.Error(t, err)
	c.URL = "\n"
	_, err = NewCallback(c, events.Signup, nil)
	assert.Error(t, err)
	c.URL = testURL
	data := types.Map{
		"test": "hello",
	}
	cb, err := NewCallback(c, events.Signup, data)
	assert.NoError(t, err)
	assert.Equal(t, events.Signup, cb.event)
	assert.Equal(t, testCallback, cb.RequestURL())
	b, err := data.JSON()
	assert.NoError(t, err)
	assert.JSONEq(t, string(b), cb.RequestBody().String())

}
