package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	mailName     = "FooBar Co"
	mailLink     = "http://mail.example.com"
	mailLogo     = "http://mail.example.com/logo.png"
	testFrom     = "foo <foo@example.com%s>"
	testTheme    = "flat"
	smtpHost     = "smtp.example.com"
	testPort     = 25
	smtpUsername = "peaches"
	smtpPassword = "secret-password!"
	smtpAuth     = "cram-md5"
	smtpEncrypt  = "tls"
	keepalive    = false
	smtpExp      = 100 * time.Minute
	smtpLimit    = 100 * time.Minute
	smtpSpam     = false
	template     = "./templates/mail.tmpl"
)

func TestMail(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		m := c.Mail
		assert.Equal(t, mailName+test.mark, m.Name)
		assert.Equal(t, mailLink+test.mark, m.Link)
		assert.Equal(t, mailLogo+test.mark, m.Logo)
		assert.Equal(t, fmt.Sprintf(testFrom, test.mark), m.From)
		assert.Equal(t, testTheme+test.mark, m.Theme)
		s := m.SMTP
		assert.Equal(t, smtpHost+test.mark, s.Host)
		assert.Equal(t, testPort, s.Port)
		assert.Equal(t, smtpUsername+test.mark, s.Username)
		assert.Equal(t, smtpPassword+test.mark, s.Password)
		assert.Equal(t, smtpAuth+test.mark, s.Authentication)
		assert.Equal(t, smtpEncrypt+test.mark, s.Encryption)
		assert.Equal(t, keepalive, s.KeepAlive)
		assert.Equal(t, smtpExp, s.Expiration)
		assert.Equal(t, smtpLimit, s.SendLimit)
		assert.Equal(t, smtpSpam, s.SpamProtection)
		assert.Equal(t, template+test.mark, m.Layout)
	})
}

// tests the ENV vars are correctly taking precedence
func TestMail_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			m := c.Mail
			assert.Equal(t, mailName, m.Name)
			assert.Equal(t, mailLink, m.Link)
			assert.Equal(t, mailLogo, m.Logo)
			assert.Equal(t, fmt.Sprintf(testFrom, ""), m.From)
			assert.Equal(t, testTheme, m.Theme)
			s := m.SMTP
			assert.Equal(t, smtpHost, s.Host)
			assert.Equal(t, testPort, s.Port)
			assert.Equal(t, smtpUsername, s.Username)
			assert.Equal(t, smtpPassword, s.Password)
			assert.Equal(t, smtpAuth, s.Authentication)
			assert.Equal(t, smtpEncrypt, s.Encryption)
			assert.Equal(t, keepalive, s.KeepAlive)
			assert.Equal(t, smtpExp, s.Expiration)
			assert.Equal(t, smtpLimit, s.SendLimit)
			assert.Equal(t, smtpSpam, s.SpamProtection)
			assert.Equal(t, template, m.Layout)
		})
	}
}

// test the *un-normalized* defaults with load
func TestMail_Defaults(t *testing.T) {
	clearEnv()
	c, err := load("")
	assert.NoError(t, err)
	def := mailDefaults
	assert.Equal(t, def, c.Mail)
}

func TestMail_Normalization(t *testing.T) {
	srv := Service{
		Name:    service,
		SiteURL: siteURL,
	}
	m := Mail{
		MailFormat: MailFormat{
			Logo: mailLogo,
		},
	}
	const (
		name = "Example"
		from = "Example <do-not-reply@example.com>"
	)
	err := m.normalize(srv)
	assert.NoError(t, err)
	assert.Equal(t, name, m.Name)
	assert.Equal(t, siteURL, m.Link)
	assert.Equal(t, mailLogo, m.Logo)
	assert.Equal(t, from, m.From)
	m.Link = "\n"
	err = m.normalize(srv)
	assert.Error(t, err)
}

func TestMail_CheckSendLimit(t *testing.T) {
	m := Mail{}
	err := m.CheckSendLimit(nil)
	assert.NoError(t, err)
	m.SendLimit = time.Millisecond
	time.Sleep(10 * time.Millisecond)
	last := time.Now()
	err = m.CheckSendLimit(&last)
	assert.Error(t, err)
	m.SendLimit = 10 * time.Minute
	err = m.CheckSendLimit(&last)
	assert.Error(t, err)
	last = last.Add(100 * time.Minute)
	err = m.CheckSendLimit(&last)
	assert.Error(t, err)
	m.SendLimit = 0
	last = time.Now()
	err = m.CheckSendLimit(&last)
	assert.NoError(t, err)
}
