package mail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"net/textproto"
	"path/filepath"
	"sync"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/utils"
	"github.com/sirupsen/logrus"
	smtp "github.com/xhit/go-simple-mail/v2"
)

// Client is a mail & smtp client
type Client struct {
	from      mail.Address
	client    *smtp.SMTPClient
	server    *smtp.SMTPServer
	config    config.Mail
	log       logrus.FieldLogger
	mu        sync.RWMutex
	keepalive *time.Ticker
	// validate email addresses against the smtp server
	smtpEmailValidation bool
}

// NewMailClient returns a new mail client.
func NewMailClient(c *config.Config, log logrus.FieldLogger) (*Client, error) {
	m := c.Mail
	if _, err := validate.Email(m.From); err != nil {
		err = fmt.Errorf("from address %w", err)
		return nil, err
	}
	// it is impossible for this to fail if validate.Email()
	// succeeds because validate.Email() calls parseAddress()
	from, _ := parseAddress(m.From)
	if log == nil {
		log = logrus.New()
	}
	if m.Host == "" {
		log = log.WithField("smtp", "offline")
		return &Client{from: from, log: log}, nil
	}
	log = log.WithField("smtp", m.Host)
	s := smtp.NewSMTPClient()
	s.Host = m.Host
	s.Port = m.Port
	s.Username = m.Username
	s.Password = m.Password
	switch m.Authentication {
	case "login":
		s.Authentication = smtp.AuthLogin
	case "cram-md5":
		s.Authentication = smtp.AuthCRAMMD5
	default:
		s.Authentication = smtp.AuthPlain
	}
	const (
		defaultPort = 25  // SMTP
		sslPort     = 465 // SMTPS (deprecated)
		tlsPort     = 587 // SMTP tls
	)
	switch m.Encryption {
	case "ssl":
		s.Encryption = smtp.EncryptionSSL
	case "tls":
		if s.Port == defaultPort || s.Port == sslPort {
			log.Warnf("using smtp tls port: %d", tlsPort)
			s.Port = tlsPort
		}
		s.Encryption = smtp.EncryptionTLS
	default:
		s.Encryption = smtp.EncryptionNone
	}
	// TODO: should we instead test the connection early, and then re-connect?
	s.KeepAlive = m.KeepAlive
	s.ConnectTimeout = 10 * time.Second
	s.SendTimeout = 10 * time.Second
	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	if c.IsDebug() {
		s.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &Client{
		from:   from,
		server: s,
		config: m,
		log:    log,
	}, nil
}

// Open opens the mail client and connects to the smtp server.
func (m *Client) Open() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	if m.client != nil {
		return nil
	}
	var err error
	m.client, err = m.server.Connect()
	if err != nil {
		return err
	}
	m.smtpEmailValidation = true
	if !m.client.KeepAlive {
		return m.close()
	}
	return m.keepAlive()
}

func (m *Client) keepAlive() error {
	if m.client == nil || !m.client.KeepAlive {
		return nil
	}
	logErr := func(err error) error {
		if err == nil {
			return nil
		}
		err = fmt.Errorf("keepalive: %w", err)
		m.log.Error(err)
		return err
	}
	sendNoop := func() error {
		err := m.client.Noop()
		if err == nil {
			return nil
		}
		e, ok := err.(*textproto.Error)
		if !ok {
			return logErr(err)
		}
		if e.Code != http.StatusOK {
			return logErr(err)
		}
		return nil
	}
	err := sendNoop()
	if err != nil {
		return nil
	}
	m.keepalive = time.NewTicker(30 * time.Second)
	go func() {
		for range m.keepalive.C {
			err = sendNoop()
			if err != nil {
				return
			}
		}
	}()
	return nil
}

// Close closes the mail client and disconnects from the smtp server.
func (m *Client) Close() error {
	return m.close()
}

// close is non-locking close and can be called from Client.Open()
func (m *Client) close() error {
	defer func() {
		m.smtpEmailValidation = false
		if m.keepalive != nil {
			m.keepalive.Stop()
		}
		m.client = nil
	}()
	if m.client == nil {
		return nil
	}
	return m.client.Close()
}

// IsOffline returns true if the smtp server is offline.
func (m *Client) IsOffline() bool {
	return m.server == nil
}

// UseSpamProtection use live smtp server email validation.
func (m *Client) UseSpamProtection(enable bool) {
	m.config.SpamProtection = enable
}

func (m *Client) validateEmailAccount() bool {
	return m.config.SpamProtection && m.smtpEmailValidation
}

// ValidateEmail validates email accounts with a live smtp server or offline.
func (m *Client) ValidateEmail(e string) (string, error) {
	if m.IsOffline() || !m.validateEmailAccount() {
		return validate.Email(e)
	}
	return validate.EmailAccount(m.server.Host, m.from.Address, e)
}

// Status returns the status of the mail client.
func (m *Client) Status() string {
	if m.IsOffline() {
		return "offline"
	} else if m.server.KeepAlive {
		return "keepalive"
	} else {
		return "idle"
	}
}

// From the from: address to use.
func (m *Client) From() string {
	return m.from.String()
}

// Send can be used to send one-off emails to users
func (m *Client) Send(to, logo, subject, html, plain string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	if _, err := m.ValidateEmail(to); err != nil {
		err = fmt.Errorf("to email %w", err)
		return err
	}
	if html == "" && plain == "" {
		return errors.New("body required")
	}
	msg := smtp.NewMSG()
	msg.SetFrom(m.From())
	msg.AddTo(to)
	msg.SetSubject(subject)
	if html != "" {
		msg.SetBody(smtp.TextHTML, html)
	}
	if plain != "" {
		if html != "" {
			msg.AddAlternative(smtp.TextPlain, plain)
		} else {
			msg.SetBody(smtp.TextPlain, plain)
		}
	}
	if utils.IsLocalPath(logo) {
		msg.AddInline(logo, filepath.Base(logo))
	}
	if m.client != nil {
		return msg.Send(m.client)
	}
	client, err := m.server.Connect()
	if err != nil {
		return err
	}
	m.smtpEmailValidation = true
	return msg.Send(client)
}

// Send can be used to send one-off templated emails to users.
func (m *Client) sendEmail(t template.Template) error {
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	err := template.LoadTemplate(m.config, t)
	if err != nil {
		return err
	}
	html, err := t.HTML()
	if err != nil {
		return err
	}
	plain, err := t.PlainText()
	if err != nil {
		return err
	}
	return m.Send(t.To(), t.Logo(), t.Subject(), html, plain)
}

// SendChangeEmail sends an email change request
func (m *Client) SendChangeEmail(to, newAddress, token, referrerURL string) error {
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	if token == "" {
		return errors.New("invalid token")
	}
	toAddr, err := parseAddress(to)
	if err != nil {
		return err
	}
	e := template.NewChangeEmail(m.config.ConfirmUser, toAddr, newAddress, token, referrerURL)
	return m.sendEmail(e)
}

// SendConfirmUser sends a mail confirmation.
func (m *Client) SendConfirmUser(to, token, referrerURL string) error {
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	if token == "" {
		return errors.New("invalid token")
	}
	toAddr, err := parseAddress(to)
	if err != nil {
		return err
	}
	e := template.NewConfirmUser(m.config.ConfirmUser, toAddr, token, referrerURL)
	return m.sendEmail(e)
}

// SendResetPassword sends an invite mail to a new user
func (m *Client) SendResetPassword(to, token, referrerURL string) error {
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	if token == "" {
		return errors.New("invalid token")
	}
	toAddr, err := parseAddress(to)
	if err != nil {
		return err
	}
	e := template.NewResetPassword(m.config.ResetPassword, toAddr, token, referrerURL)
	return m.sendEmail(e)
}

// SendInviteUser sends an invite mail to a new user
func (m *Client) SendInviteUser(from, to, token, referrerURL string) error {
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	if token == "" {
		return errors.New("invalid token")
	}
	fromAddr, err := optionalAddress(from)
	if err != nil {
		return err
	}
	toAddr, err := parseAddress(to)
	if err != nil {
		return err
	}
	e := template.NewInviteUser(m.config.InviteUser, fromAddr, toAddr, token, referrerURL)
	return m.sendEmail(e)
}

// SendSignupCode sends an invite mail to a new user
func (m *Client) SendSignupCode(from, to, token, referrerURL string) error {
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	if token == "" {
		return errors.New("invalid signup code")
	}
	fromAddr, err := optionalAddress(from)
	if err != nil {
		return err
	}
	toAddr, err := parseAddress(to)
	if err != nil {
		return err
	}
	e := template.NewSignupCode(m.config.SignupCode, fromAddr, toAddr, token, referrerURL)
	return m.sendEmail(e)
}

func parseAddress(address string) (mail.Address, error) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return mail.Address{}, err
	}
	return *addr, err
}

func optionalAddress(address string) (mail.Address, error) {
	if address == "" {
		return mail.Address{}, nil
	}
	return parseAddress(address)
}
