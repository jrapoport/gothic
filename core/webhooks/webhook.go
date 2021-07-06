package webhooks

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/store/types"
)

// CallWebhook calls a webhook callback with an event and msg.
func CallWebhook(c config.Webhooks, evt events.Event, msg types.Map) error {
	cb, err := NewCallback(c, evt, msg)
	if err != nil {
		return err
	}
	req, err := NewWebhookRequest(cb)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: c.Timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body := strings.TrimRight(string(b), "\n")
		err = fmt.Errorf("%s: %s", res.Status, body)
		return err
	}
	return nil
}
