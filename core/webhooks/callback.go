package webhooks

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/jwt"
	"github.com/jrapoport/gothic/models/types"
)

// Callback holds a webhook callback.
type Callback struct {
	url     *url.URL
	claims  *jwt.WebhookClaims
	event   events.Event
	jwt     config.JWT
	payload []byte
}

// NewCallback returns a new webhook callback
func NewCallback(c config.Webhooks, e events.Event, msg types.Map) (*Callback, error) {
	if e == events.Unknown || e == events.All {
		return nil, errors.New("invalid event")
	}
	if c.URL == "" {
		return nil, errors.New("invalid url")
	}
	u, err := url.Parse(c.URL)
	if err != nil {
		return nil, err
	}
	payload := msg.JSON()
	sum := checksum(payload)
	claims := jwt.NewWebhookClaims(sum)
	return &Callback{
		event:   e,
		url:     u,
		claims:  claims,
		jwt:     c.JWT,
		payload: payload,
	}, nil
}

// RequestClaims returns the claims for the webhook callback.
func (c Callback) RequestClaims() jwt.Claims {
	return c.claims
}

// RequestURL the request url for the webhook callback
func (c Callback) RequestURL() string {
	return config.FormatWebhookURL(c.url.String(), c.event)
}

// RequestBody the request body for the webhook callback
func (c Callback) RequestBody() *bytes.Buffer {
	return bytes.NewBuffer(c.payload)
}

func checksum(data []byte) string {
	sha := sha256.New()
	sha.Write(data)
	return hex.EncodeToString(sha.Sum(nil))
}
