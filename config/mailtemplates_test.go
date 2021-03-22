package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	linkFormat  = "/:action/:token/link"
	subject     = "Email Subject"
	referralURL = "http://referral.example.com"
)

func TestMailTemplates(t *testing.T) {
	runTests(t, func(t *testing.T, test testCase, c *Config) {
		m := c.Mail
		tmpls := []MailTemplate{
			m.ChangeEmail,
			m.ConfirmUser,
			m.InviteUser,
			m.SignupCode,
			m.ResetPassword,
		}
		for _, tmpl := range tmpls {
			assert.Equal(t, template+test.mark, tmpl.Template)
			assert.Equal(t, linkFormat+test.mark, tmpl.LinkFormat)
			assert.Equal(t, subject+test.mark, tmpl.Subject)
			assert.Equal(t, referralURL+test.mark, tmpl.ReferralURL)
		}
	})
}

// tests the ENV vars are correctly taking precedence
func TestMailTemplates_Env(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clearEnv()
			loadDotEnv(t)
			c, err := loadNormalized(test.file)
			assert.NoError(t, err)
			m := c.Mail
			tmpls := []MailTemplate{
				m.ChangeEmail,
				m.ConfirmUser,
				m.InviteUser,
				m.SignupCode,
				m.ResetPassword,
			}
			for _, tmpl := range tmpls {
				assert.Equal(t, template, tmpl.Template)
				assert.Equal(t, linkFormat, tmpl.LinkFormat)
				assert.Equal(t, subject, tmpl.Subject)
				assert.Equal(t, referralURL, tmpl.ReferralURL)
			}
		})
	}
}

func TestMailTemplates_Normalization(t *testing.T) {
	srv := Service{
		SiteURL: siteURL,
	}
	mt := MailTemplates{}
	err := mt.normalize(srv)
	assert.NoError(t, err)
	referralURLs := []string{
		mt.ChangeEmail.ReferralURL,
		mt.ConfirmUser.ReferralURL,
		mt.InviteUser.ReferralURL,
		mt.ResetPassword.ReferralURL,
		mt.SignupCode.ReferralURL,
	}
	for _, ref := range referralURLs {
		assert.Equal(t, siteURL, ref)
	}
	mt.ChangeEmail.ReferralURL = "\n"
	err = mt.normalize(srv)
	assert.Error(t, err)
}

func TestFormatLink(t *testing.T) {
	const action = "the-action"
	const token = "1234567890asdfghjklqwertyuiopzxcvbnm="
	var tmplTests = []struct {
		link     string
		expected string
	}{
		{"", ""},
		{
			"http://www.example.com",
			"http://www.example.com",
		},
		{
			"http://www.example.com/:action?q=go+language#foo%26bar",
			"http://www.example.com/" + action + "?q=go+language#foo%26bar",
		},
		{
			"http://www.example.com/?q=go+language#foo%26bar:token",
			"http://www.example.com/?q=go+language#foo%26bar" + token,
		},
		{
			"http://www.example.com/:action/go?q=go+language#foo%26bar:token",
			"http://www.example.com/" + action + "/go?q=go+language#foo%26bar" + token,
		},
		{
			"http://www.example.com/:action/go?q=go+language&token=:token#foo%26bar",
			"http://www.example.com/" + action + "/go?q=go+language&token=" + token + "#foo%26bar",
		},
		{
			"a/b/:action/c/:token/",
			"a/b/" + action + "/c/" + token + "/",
		},
		{
			"http://[fe80::1%25%65%6e%301-._~]:8080/:action/?:token",
			"http://[fe80::1%25%65%6e%301-._~]:8080/" + action + "/?" + token,
		},
		{
			"scheme://!$&'()*+,;=hello!:1/:action_path#:token",
			"scheme://!$&'()*+,;=hello!:1/" + action + "_path#" + token,
		},
		{
			"http://host/!$&'()*+,:token;=:@[:action]",
			"http://host/!$&'()*+," + token + ";=:@[" + action + "]",
		},
		{
			"http://[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:8080/foo:action:token",
			"http://[2b01:e34:ef40:7730:8e70:5aff:fefe:edac]:8080/foo" + action + token,
		},
		{
			"http://hello.世界.com/foo/:Action:action?:Token=:token&referral=localhost",
			"http://hello.世界.com/foo/:Action" + action + "?:Token=" + token + "&referral=localhost",
		},
		{
			"mailto:?subject=:action&body=here%27s%20the%20link%20you%20requested%3A%20:token",
			"mailto:?subject=" + action + "&body=here%27s%20the%20link%20you%20requested%3A%20" + token,
		},
	}
	for i, test := range tmplTests {
		l := FormatLink(test.link, action, token)
		assert.Equal(t, test.expected, l, i)
	}
	clearTests := []struct {
		action   string
		token    string
		expected string
	}{
		{"", "", "http://www.example.com//#?foo=bar"},
		{action, "", "http://www.example.com/the-action/#?foo=bar"},
		{"", token, "http://www.example.com//#1234567890asdfghjklqwertyuiopzxcvbnm=?foo=bar"},
		{action, token, "http://www.example.com/the-action/#1234567890asdfghjklqwertyuiopzxcvbnm=?foo=bar"},
	}
	const noLink = "http://www.example.com/:action/#:token?foo=bar"
	for _, test := range clearTests {
		l := FormatLink(noLink, test.action, test.token)
		assert.Equal(t, test.expected, l)
	}
}
