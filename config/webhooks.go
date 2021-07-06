package config

import (
	"net/url"
	"strings"
	"time"

	"github.com/jrapoport/gothic/core/events"
)

// Webhooks config
type Webhooks struct {
	Events []events.Event `json:"events"`
	// URL is the url for the webhook callback. If the url includes an ':event'
	// string it will automatically be replaced with name of the callback event
	URL        string        `json:"url"`
	MaxRetries int           `json:"max_retries" yaml:"max_retries" mapstructure:"max_retries"`
	Timeout    time.Duration `json:"timeout"`
	JWT        `yaml:",inline" mapstructure:",squash"`
}

func (w *Webhooks) normalize(srv Service, j JWT) (err error) {
	if w.URL == "" {
		return nil
	}
	_, err = url.Parse(w.URL)
	if err != nil {
		return err
	}
	for _, e := range w.Events {
		if e != events.All {
			continue
		}
		w.Events = []events.Event{events.All}
		break
	}
	return w.JWT.normalize(srv, j)
}

// Enabled returns true if the enabled.
func (w Webhooks) Enabled() bool {
	return w.URL != ""
}

// HasEvent returns true if the event is configured.
func (w Webhooks) HasEvent(event events.Event) bool {
	for _, name := range w.Events {
		if event == name {
			return true
		}
	}
	return false
}

// WebhookURLEvent is the webhook event token
const WebhookURLEvent = ":event"

// FormatWebhookURL formats the callback URL. If the url contains an
// ':event' string, it will be replaced with the name of the event.
func FormatWebhookURL(url string, event events.Event) string {
	return strings.ReplaceAll(url, WebhookURLEvent, string(event))
}
