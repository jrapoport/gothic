package core

import (
	"errors"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/core/webhooks"
	"github.com/jrapoport/gothic/store/types"
)

// LoadWebhooks load the webhooks.
func (a *API) LoadWebhooks() error {
	wh := a.config.Webhook
	if !wh.Enabled() {
		a.log.Warn("webhooks disabled")
		return nil
	}
	if wh.Secret == "" {
		err := errors.New("webhook secret required")
		return a.logError(err)
	}
	_, err := url.Parse(wh.URL)
	if err != nil {
		return a.logError(err)
	}
	hooks := wh.Events
	if len(hooks) <= 0 {
		a.log.Warn("webhook events not found")
		return nil
	}
	for _, e := range hooks {
		if e != events.All {
			continue
		}
		hooks = []events.Event{events.All}
		break
	}
	for _, e := range hooks {
		a.AddListener(e, a.callback)
	}
	return nil
}

func (a *API) callback(evt events.Event, msg types.Map) {
	wc := a.config.Webhook
	expBack := backoff.NewExponentialBackOff()
	max := wc.MaxRetries
	b := backoff.WithMaxRetries(expBack, uint64(max))
	err := backoff.RetryNotify(
		func() error {
			return webhooks.CallWebhook(wc, evt, msg)
		},
		b,
		func(err error, duration time.Duration) {
			a.log.
				WithError(err).
				Warn("webhook failed")
			a.log.
				WithField("duration", duration).
				Info("webhook retry in...")
		},
	)
	_ = a.logError(err)
}
