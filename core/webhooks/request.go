package webhooks

import (
	"errors"
	"net/http"

	"github.com/jrapoport/gothic/jwt"
)

// HTTP Headers
const (
	WebhookSignature = "x-webhook-signature"
	ContentType      = "Content-Type"
	JSONContent      = "application/json"
)

// NewWebhookRequest returns a new http request for the webhook callback.
func NewWebhookRequest(cb *Callback) (*http.Request, error) {
	if cb == nil {
		return nil, errors.New("invalid callback")
	}
	req, err := http.NewRequest(http.MethodPost, cb.RequestURL(), cb.RequestBody())
	if err != nil {
		return nil, err
	}
	req.Header.Set(ContentType, JSONContent)
	token := jwt.NewToken(cb.jwt, cb.RequestClaims())
	bearer, err := token.Bearer()
	if err != nil {
		return nil, err
	}
	req.Header.Set(WebhookSignature, bearer)
	return req, nil
}
