package webhooks

import (
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/tokens/jwt"
)

// WebhookClaims the jwt claims for a webhook.
type WebhookClaims struct {
	jwt.StandardClaims
	Checksum string `json:"chk"`
}

var _ jwt.Claims = (*WebhookClaims)(nil)

// NewWebhookClaims returns a new set of webhook jwt claims.
func NewWebhookClaims(c config.JWT, checksum string) WebhookClaims {
	std := jwt.NewStandardClaims(c)
	std.Subject = "webhook"
	return WebhookClaims{
		StandardClaims: std,
		Checksum:       checksum,
	}
}
