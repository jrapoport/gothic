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
	"github.com/jrapoport/gothic/log"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/utils"
	smtp "github.com/xhit/go-simple-mail/v2"
)

// Client is a mail & smtp client
type Client struct {
	from      mail.Address
	client    *smtp.SMTPClient
	server    *smtp.SMTPServer
	config    config.Mail
	log       log.Logger
	mu        sync.RWMutex
	keepalive *time.Ticker
	// validate email addresses against the smtp server
	smtpEmailValidation bool
}

// NewMailClient returns a new mail client.
func NewMailClient(c *config.Config, l log.Logger) (*Client, error) {
	m := c.Mail
	if _, err := validate.Email(m.From); err != nil {
		err = fmt.Errorf("from address %w", err)
		return nil, err
	}
	// it is impossible for this to fail if validate.Email()
	// succeeds because validate.Email() calls parseAddress()
	from, _ := parseAddress(m.From)
	if l == nil {
		l = log.New()
	}
	if m.Host == "" {
		l = l.WithName("smtp-offline")
		return &Client{from: from, log: l}, nil
	}
	l = l.WithName("smtp-" + m.Host)
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
			l.Warnf("using smtp tls port: %d", tlsPort)
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
		log:    l,
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
	m.keepalive = time.NewTicker(30 * time.Second)
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		wait.Done()
		err := m.noop()
		if err != nil {
			return
		}
		for range m.keepalive.C {
			err = m.noop()
			if err != nil {
				return
			}
		}
	}()
	wait.Wait()
	return nil
}

func (m *Client) noop() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.client == nil {
		return nil
	}
	err := m.client.Noop()
	if err == nil {
		return nil
	}
	e, ok := err.(*textproto.Error)
	if ok && e.Code == http.StatusOK {
		return nil
	}
	err = fmt.Errorf("keepalive: %w", err)
	m.log.Error(err)
	return err
}

// Close closes the mail client and disconnects from the smtp server.
func (m *Client) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.close()
}

// close is non-locking close and can be called from Client.Open()
func (m *Client) close() error {
	m.smtpEmailValidation = false
	if m.keepalive != nil {
		m.keepalive.Stop()
	}
	if m.client == nil {
		return nil
	}
	err := m.client.Close()
	m.client = nil
	return err
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
	m.mu.Lock()
	defer m.mu.Unlock()
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
	if subject == "" {
		subject = m.defaultSubject()
	}
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
func (m *Client) sendTemplate(t template.Template) error {
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

type Type int

const (
	HTML Type = iota
	Markdown
	Template
)

type Content struct {
	Type      Type
	Body      string
	Plaintext string
}

func (m *Client) SendEmail(to, subject string, content Content) error {
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	if content.Body == "" {
		return errors.New("content body required")
	}
	var err error
	switch content.Type {
	case Template:
		err = m.SendTemplateEmail(to, subject, content.Body)
	case Markdown:
		err = m.SendMarkdownEmail(to, subject, content.Body)
	case HTML:
		fallthrough
	default:
		err = m.Send(to, m.config.Logo, subject, content.Body, content.Plaintext)
	}
	return err
}

// SendTemplateEmail sends an generic email based on a body template
func (m *Client) SendTemplateEmail(to, subject, body string) error {
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	toAddr, err := parseAddress(to)
	if err != nil {
		return err
	}
	if body == "" {
		return errors.New("body template required")
	}
	c := config.MailTemplate{
		Subject: subject,
	}
	e := template.NewEmail(c, toAddr, body)
	return m.sendTemplate(e)
}

// SendMarkdownEmail sends an markdown email
func (m *Client) SendMarkdownEmail(to, subject, markdown string) error {
	if m.IsOffline() {
		m.log.Warn("mail client is offline")
		return nil
	}
	toAddr, err := parseAddress(to)
	if err != nil {
		return err
	}
	if markdown == "" {
		return errors.New("email markdown required")
	}
	c := config.MailTemplate{
		Subject: subject,
	}
	e := template.NewMarkdownMail(c, toAddr, markdown)
	return m.sendTemplate(e)
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
	return m.sendTemplate(e)
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
	return m.sendTemplate(e)
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
	return m.sendTemplate(e)
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
	return m.sendTemplate(e)
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
	return m.sendTemplate(e)
}

func (m *Client) defaultSubject() string {
	name := "us"
	if m.config.Name != "" {
		name = m.config.Name
	}
	return "An important message from " + name
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
