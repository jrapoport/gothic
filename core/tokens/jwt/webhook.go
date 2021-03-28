package jwt

const ChecksumKey = "chk"

// WebhookClaims the jwt claims for a webhook.
type WebhookClaims struct {
	StandardClaims
	Checksum string `json:"chk"`
}

var _ Claims = (*WebhookClaims)(nil)

// NewWebhookClaims returns a new set of webhook jwt claims.
func NewWebhookClaims(checksum string) *WebhookClaims {
	std := *NewStandardClaims("webhook")
	return &WebhookClaims{
		StandardClaims: std,
		Checksum:       checksum,
	}
}

func (c WebhookClaims) ParseToken(tok *Token) {
	c.StandardClaims.ParseToken(tok)
	if v, ok := c.Get(ChecksumKey); ok {
		c.Checksum = v.(string)
	}
}
