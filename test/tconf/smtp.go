package tconf

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/log"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	smtpHost = "localhost"
	smtpPort = 2525
)

// SMTPMock is a mock smtp server.
type SMTPMock struct {
	smtp  guerrilla.Daemon
	hooks sync.Map
}

// AddHook adds a mock hook.
func (m *SMTPMock) AddHook(t *testing.T, hook func(email string)) {
	key := utils.RandomUsername()
	m.hooks.Store(key, hook)
	t.Cleanup(func() {
		m.hooks.LoadAndDelete(key)
	})
}

var mu sync.Mutex

// MockSMTP returns a mocked smtp instance.
func MockSMTP(t *testing.T, c *config.Config) (*config.Config, *SMTPMock) {
	c.Mail.Host = smtpHost
	c.Mail.Port = smtpPort
	cfg := &guerrilla.AppConfig{
		BackendConfig: backends.BackendConfig{
			"save_process": "mock|Debugger",
		},
		AllowedHosts: []string{
			".",
		},
		LogLevel: c.Logger.Level,
	}
	const maxTries = 100
	var port int
	mu.Lock()
	for {
		addr := fmt.Sprintf("%s:%d",
			c.Mail.Host, c.Mail.Port+port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			port++
			if port > maxTries {
				t.Fatalf("cannot acquire smtp test port after %d tries", maxTries)
			}
			continue
		}
		_ = lis.Close()
		c.Mail.Port += port
		sc := guerrilla.ServerConfig{
			ListenInterface: addr,
			IsEnabled:       true,
		}
		cfg.Servers = append(cfg.Servers, sc)
		break
	}
	mu.Unlock()
	hl := &log.HookedLogger{Logger: (c.Log().(*logrus.Entry)).Logger}
	smtp := guerrilla.Daemon{
		Config: cfg,
		Logger: hl,
	}
	mock := &SMTPMock{
		smtp:  smtp,
		hooks: sync.Map{},
	}
	smtp.AddProcessor("mock", func() backends.Decorator {
		return func(p backends.Processor) backends.Processor {
			return backends.ProcessWith(
				func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {
					if task == backends.TaskSaveMail {
						mock.hooks.Range(func(_, value interface{}) bool {
							hook := value.(func(string))
							hook(e.Data.String())
							return true
						})
					}
					return p.Process(e, task)
				})
		}
	})
	err := smtp.Start()
	require.NoError(t, err)
	t.Cleanup(func() {
		smtp.Shutdown()
	})
	return c, mock
}

// GetEmailToken parses a token from an email link.
func GetEmailToken(action, email string) string {
	rx := regexp.MustCompile(`(?m)` + action + `\/([a-zA-Z0-9_\-.=\n]+)`)
	matches := rx.FindStringSubmatch(email)
	if len(matches) < 2 {
		return ""
	}
	tok := matches[len(matches)-1]
	tok = strings.ReplaceAll(tok, "=\n", "")
	return tok
}
