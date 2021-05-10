package tconf

import (
	"net"
	"regexp"
	"strconv"
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

const smtpHost = "127.0.0.1"

// SMTPMock is a mock smtp server.
type SMTPMock struct {
	smtp  guerrilla.Daemon
	hooks sync.Map
	mu    sync.Mutex
}

// AddHook adds a mock hook.
func (m *SMTPMock) AddHook(t *testing.T, hook func(email string)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := utils.RandomUsername()
	m.hooks.Store(key, hook)
	t.Cleanup(func() {
		m.hooks.LoadAndDelete(key)
	})
}

var mu sync.Mutex

// MockSMTP returns a mocked smtp instance.
func MockSMTP(t *testing.T, c *config.Config) (*config.Config, *SMTPMock) {
	mu.Lock()
	defer mu.Unlock()
	c.Mail.Host = smtpHost
	cfg := &guerrilla.AppConfig{
		BackendConfig: backends.BackendConfig{
			"save_process": "mock|Debugger",
		},
		AllowedHosts: []string{
			".",
		},
		LogLevel: c.Logger.Level,
	}
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := lis.Addr().String()
	_ = lis.Close()
	parts := strings.Split(addr, ":")
	c.Mail.Port, err = strconv.Atoi(parts[1])
	require.NoError(t, err)
	sc := guerrilla.ServerConfig{
		ListenInterface: addr,
		IsEnabled:       true,
	}
	cfg.Servers = append(cfg.Servers, sc)
	l := logrus.New()
	level, err := logrus.ParseLevel(c.Level)
	require.NoError(t, err)
	l.SetLevel(level)
	smtp := guerrilla.Daemon{
		Config: cfg,
		Logger: &log.HookedLogger{Logger: l},
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
						mock.mu.Lock()
						defer mock.mu.Unlock()
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
	err = smtp.Start()
	require.NoError(t, err)
	t.Cleanup(func() {
		smtp.Shutdown()
	})
	return c, mock
}

// GetEmailToken parses a token from an email link.
func GetEmailToken(action, email string) string {
	rx := regexp.MustCompile(`(?m)` + action + `\/#\/([a-zA-Z0-9_\-.=\n]+)`)
	matches := rx.FindStringSubmatch(email)
	if len(matches) < 2 {
		return ""
	}
	tok := matches[len(matches)-1]
	tok = strings.ReplaceAll(tok, "#", "")
	tok = strings.ReplaceAll(tok, "=\n", "")
	return tok
}
