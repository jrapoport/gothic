package jwt

const ChecksumKey = "chk"

// WebhookClaims the jwt claims for a webhook.
type WebhookClaims struct {
	StandardClaims
}

var _ Claims = (*WebhookClaims)(nil)

// NewWebhookClaims returns a new set of webhook jwt claims.
func NewWebhookClaims(checksum string) *WebhookClaims {
	c := &WebhookClaims{
		StandardClaims: *NewStandardClaims("webhook"),
	}
	_ = c.Set(ChecksumKey, checksum)
	return c
}

// Checksum returns the checksum
func (c WebhookClaims) Checksum() string {
	checksum, ok := c.Get(ChecksumKey)
	if !ok {
		return ""
	}
	return checksum.(string)
}
