package mail

import (
	"net/mail"
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	logoFile = "template/testdata/template_logo.png"

	toUsername   = "the_real_mr_flibble"
	toEmail      = "mr.flibble@example.com"
	fromUsername = "el peaches"
	fromEmail    = "el_peaches@example.com"
	testToken    = "1234567890asdfghjklqwertyuiopzxcvbnm="
	testReferral = "http://test.example.com:3000/"
)

var (
	emptyAddress = mail.Address{}
	toAddress    = mail.Address{Name: toUsername, Address: toEmail}
	fromAddress  = mail.Address{Name: fromUsername, Address: fromEmail}
)

type testCase struct {
	to   mail.Address
	from mail.Address
	tok  string
	ref  string
	Err  assert.ErrorAssertionFunc
}

func mockSMTP(t *testing.T) (*config.Config, *tconf.SMTPMock) {
	c := tconf.Config(t)
	c.Mail.Logo = logoFile
	return tconf.MockSMTP(t, c)
}

func testMailer(t *testing.T, c *config.Config) (*Client, error) {
	c.Mail.SpamProtection = false
	m, err := NewMailClient(c, c.Log())
	if err != nil {
		return nil, err
	}
	err = m.Open()
	assert.NoError(t, err)
	t.Cleanup(func() {
		err = m.Close()
		assert.NoError(t, err)
	})
	return m, nil
}

type ClientTestSuite struct {
	suite.Suite
	offline   bool
	keepalive bool
	client    *Client
	mock      *tconf.SMTPMock
	c         *config.Config
	reset     config.Config
}

func TestMailer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		offline   bool
		keepalive bool
	}{
		{"Offline", true, false},
		{"Idle", false, false},
		{"Keepalive", false, true},
	}
	for _, test := range tests {
		ts := &ClientTestSuite{
			offline:   test.offline,
			keepalive: test.keepalive,
		}
		t.Run(test.name, func(t *testing.T) {
			if !test.keepalive {
				t.Parallel()
			}
			suite.Run(t, ts)
		})
	}
}

func (ts *ClientTestSuite) SetupSuite() {
	var client *Client
	var err error
	ts.c, ts.mock = mockSMTP(ts.T())
	ts.c.Mail.KeepAlive = false
	if ts.offline {
		ts.c.Mail.Host = ""
		client, err = testMailer(ts.T(), ts.c)
	} else {
		ts.c.Mail.KeepAlive = ts.keepalive
		client, err = testMailer(ts.T(), ts.c)
	}
	ts.Require().NoError(err)
	ts.Require().NotNil(client)
	ts.client = client
}

func (ts *ClientTestSuite) TearDownSuite() {
	if ts.client == nil {
		return
	}
	err := ts.client.Close()
	ts.NoError(err)
}

func (ts *ClientTestSuite) TestNew() {
	tests := []struct {
		auth      string
		encrypt   string
		keepalive bool
		Err       assert.ErrorAssertionFunc
	}{
		{"plain", "none", false, assert.NoError},
		{"plain", "tls", false, assert.NoError},
		{"plain", "ssl", false, assert.Error},
		{"plain", "none", true, assert.NoError},
		{"plain", "tls", true, assert.NoError},
		{"plain", "ssl", true, assert.Error},
		{"login", "none", false, assert.NoError},
		{"login", "tls", false, assert.NoError},
		{"login", "ssl", false, assert.Error},
		{"login", "none", true, assert.NoError},
		{"login", "tls", true, assert.NoError},
		{"login", "ssl", true, assert.Error},
		{"cram-md5", "none", false, assert.NoError},
		{"cram-md5", "tls", false, assert.NoError},
		{"cram-md5", "ssl", false, assert.Error},
		{"cram-md5", "none", true, assert.NoError},
		{"cram-md5", "tls", true, assert.NoError},
		{"cram-md5", "ssl", true, assert.Error},
	}
	for _, test := range tests {
		c := *ts.c
		c.Mail.Authentication = test.auth
		c.Mail.Encryption = test.encrypt
		m, err := NewMailClient(&c, nil)
		ts.NoError(err)
		if err == nil {
			err = m.Open()
			if m.IsOffline() {
				ts.NoError(err)
			} else {
				test.Err(ts.T(), err)
			}
			err = m.Close()
			ts.NoError(err)
		}
		c.Mail.From = ""
		_, err = NewMailClient(&c, nil)
		ts.Error(err)
	}
	m, err := testMailer(ts.T(), ts.c)
	ts.NoError(err)
	err = m.Open()
	ts.NoError(err)
	c := *ts.c
	c.Mail.Encryption = "tls"
	c.Mail.Port = 25
	m, err = NewMailClient(&c, nil)
	ts.NoError(err)
	err = m.Open()
	if m.IsOffline() {
		ts.NoError(err)
	} else {
		ts.Error(err)
	}
}

func (ts *ClientTestSuite) TestValidateEmailAccount() {
	tests := []struct {
		email string
		Err   assert.ErrorAssertionFunc
	}{
		{"", assert.Error},
		{"@", assert.Error},
		{fromAddress.Address, assert.NoError},
		{fromAddress.Name, assert.Error},
		{fromAddress.String(), assert.NoError},
	}
	m, err := testMailer(ts.T(), ts.c)
	ts.NoError(err)
	ts.NotNil(m)
	for _, test := range tests {
		m.UseSpamProtection(false)
		m.smtpEmailValidation = false
		_, err = m.ValidateEmail(test.email)
		test.Err(ts.T(), err)
	}
	for _, test := range tests {
		// currently these will all fail.
		m.UseSpamProtection(true)
		m.smtpEmailValidation = true
		_, err = m.ValidateEmail(test.email)
		if m.IsOffline() {
			test.Err(ts.T(), err)
		} else {
			ts.Error(err)
		}
	}
}

func (ts *ClientTestSuite) TestSend() {
	const (
		testSubject = "Test Subject"
		testHTML    = "<html></html>"
		testPlain   = "--plain test--"
		testLogo    = "http://example.com/logo.png"
	)
	from := "admin@example.com"
	to := tutils.RandomEmail()
	tests := []struct {
		from  string
		to    string
		sub   string
		html  string
		plain string
		Err   assert.ErrorAssertionFunc
	}{
		{"", "", "", "", "", assert.Error},
		{from, "", "", "", "", assert.Error},
		{from, to, "", "", "", assert.Error},
		{from, to, testSubject, "", "", assert.Error},
		{from, to, testSubject, testHTML, "", assert.NoError},
		{from, to, testSubject, "", testPlain, assert.NoError},
		{from, to, testSubject, testHTML, testPlain, assert.NoError},
		{from, to, "", testHTML, testPlain, assert.NoError},
	}
	for _, test := range tests {
		err := ts.client.Send(test.to, testLogo, test.sub, test.html, test.plain)
		if ts.client.IsOffline() {
			ts.NoError(err)
		} else {
			test.Err(ts.T(), err)
		}
	}
}

func (ts *ClientTestSuite) TestSendChangeEmail() {
	ts.testSendChangeEmail()
}

func (ts *ClientTestSuite) TestSendConfirmUser() {
	ts.testSendConfirmUser()
}

func (ts *ClientTestSuite) TestSendInviteUser() {
	ts.testSendInviteUser()
}

func (ts *ClientTestSuite) TestSendResetPassword() {
	ts.testSendResetPassword()
}

func (ts *ClientTestSuite) TestSendSignupCode() {
	ts.testSendSignupCode()
}

func (ts *ClientTestSuite) TestKeepalive() {
	if !ts.keepalive {
		return
	}
	type testFunc func()
	tests := []struct {
		name string
		fn   testFunc
	}{
		{"SendChangeEmail", ts.testSendChangeEmail},
		{"SendConfirmUser", ts.testSendConfirmUser},
		{"SendInviteUser", ts.testSendInviteUser},
		{"SendResetPassword", ts.testSendResetPassword},
		{"SendSignupCode", ts.testSendSignupCode},
	}
	ts.client.keepalive.Reset(1 * time.Second)
	for _, test := range tests {
		ts.Run(test.name, func() {
			time.Sleep(500 * time.Millisecond)
			test.fn()
		})
	}
}

func (ts *ClientTestSuite) sendTest(send func(tc testCase) error) {
	var sendTests = []testCase{
		{toAddress, emptyAddress, "", testReferral, assert.Error},
		{toAddress, fromAddress, "", testReferral, assert.Error},
		{toAddress, fromAddress, testToken, testReferral, assert.NoError},
	}
	for _, test := range sendTests {
		var sent string
		ts.mock.AddHook(ts.T(), func(email string) {
			sent = email
		})
		err := send(test)
		if ts.client.IsOffline() {
			ts.NoError(err)
		} else {
			test.Err(ts.T(), err)
			if err == nil && test.tok != "" {
				ts.Eventually(func() bool {
					return sent != ""
				}, 200*time.Millisecond, 10*time.Millisecond)
				ts.Contains(sent, test.ref)
				ts.Contains(sent, test.tok)

			}
		}
		test.to.Address = ""
		err = send(test)
		if ts.client.IsOffline() {
			ts.NoError(err)
		} else {
			ts.Error(err)
		}
	}
}

func (ts *ClientTestSuite) testSendChangeEmail() {
	ts.sendTest(func(tc testCase) error {
		return ts.client.SendChangeEmail(tc.to.String(), tc.to.String(), tc.tok, tc.ref)
	})
}

func (ts *ClientTestSuite) testSendConfirmUser() {
	ts.sendTest(func(tc testCase) error {
		return ts.client.SendConfirmUser(tc.to.String(), tc.tok, tc.ref)
	})
}

func (ts *ClientTestSuite) testSendInviteUser() {
	ts.sendTest(func(tc testCase) error {
		return ts.client.SendInviteUser(tc.from.String(), tc.to.String(), tc.tok, tc.ref)
	})
}

func (ts *ClientTestSuite) testSendResetPassword() {
	ts.sendTest(func(tc testCase) error {
		return ts.client.SendResetPassword(tc.to.String(), tc.tok, tc.ref)
	})
}

func (ts *ClientTestSuite) testSendSignupCode() {
	ts.sendTest(func(tc testCase) error {
		return ts.client.SendSignupCode(tc.from.String(), tc.to.String(), tc.tok, tc.ref)
	})
}
