package conf

import "github.com/dgrijalva/jwt-go/v4"

type WebhookConfig struct {
	URL        string   `json:"url"`
	Retries    int      `json:"retries"`
	TimeoutSec int      `json:"timeout_sec"`
	Secret     string   `json:"secret"`
	Method     string   `json:"method" default:"HS256"`
	Events     []string `json:"events"`
}

func (w *WebhookConfig) SigningMethod() jwt.SigningMethod {
	m := "HS256"
	if w.Method != "" {
		m = w.Method
	}
	return jwt.GetSigningMethod(m)
}

func (w *WebhookConfig) HasEvent(event string) bool {
	for _, name := range w.Events {
		if event == name {
			return true
		}
	}
	return false
}
