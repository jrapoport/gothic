package template

import (
	"fmt"
	"io/ioutil"
	"net/mail"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDir    = "testdata/"
	testBody   = testDir + "template_body.yaml"
	testLayout = testDir + "template_layout.yaml"
	testLogo   = testDir + "template_logo.png"
)

const (
	confSub  = "test subject"
	confLink = "https://example.com"
	confName = "Test Service"
	// yaml file values
	tmplSub  = ""
	tmplLink = "https://template.example.com"
	tmplName = "template example"
	tmplFrom = "admin@example.com"
	tmplLogo = "https://template.example.com/logo.png"
)

const (
	testAction = "the-action"
	testTok    = "1234567890asdfghjklqwertyuiopzxcvbnm="
	testRef    = "https://test.example.com:3000/"
	linkFmt    = "/:action/#token=:token"
	testLink2  = "https://test.example.com:3000/the-action/#token=1234567890asdfghjklqwertyuiopzxcvbnm="
)

const (
	toName    = "the_real_mr_flibble"
	toEmail   = "foo@example.com"
	fromName  = "El Peaches"
	fromEmail = "bar@example.com"
)

var (
	emptyAddr = mail.Address{}
	nameAddr  = mail.Address{Name: toName}
	emailAddr = mail.Address{Address: toEmail}
	toAddr    = mail.Address{Name: toName, Address: toEmail}
	fromAddr  = mail.Address{Name: fromName, Address: fromEmail}
)

type testCase struct {
	mc   config.Mail
	tmpl string
	to   mail.Address
	from mail.Address
	tok  string
	ref  string
	Err  assert.ErrorAssertionFunc
}

func assertEqualContent(t *testing.T, path string, content string, msgAndArgs ...interface{}) bool {
	// NOTE: uncomment this line to regenerate the test files.
	// _ = ioutil.WriteFile(path, []byte(content), 0600)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		msgFormat := fmt.Sprintf("error when running ioutil.ReadFile(%q): %s", path, err)
		return assert.Fail(t, msgFormat, msgAndArgs...)
	}
	expected := string(data)
	// this is a kind of a gross hack but w/e
	rex := regexp.MustCompile(`Copyright.*Gothic`)
	content = rex.ReplaceAllString(content, "")
	expected = rex.ReplaceAllString(expected, "")
	return assert.Equal(t, content, expected, msgAndArgs)
}

func testTemplateLoad(t *testing.T, newTemplate func(sub string, test testCase) Template) {
	c := tconf.Config(t).Mail
	const dne = "does-not-exist" //
	tests := []testCase{
		{c, "", emailAddr, emptyAddr, testTok, testRef, assert.NoError},
		{c, "", toAddr, emptyAddr, testTok, "", assert.NoError},
		{c, "", toAddr, emptyAddr, "", testRef, assert.NoError},
		{c, "", toAddr, emptyAddr, testTok, testRef, assert.NoError},
		{c, "", toAddr, fromAddr, testTok, testRef, assert.NoError},
		{c, "", toAddr, fromAddr, testTok, "", assert.NoError},
		{c, "", toAddr, fromAddr, "", testRef, assert.NoError},
		{c, "", toAddr, fromAddr, testTok, testRef, assert.NoError},
		{c, dne, toAddr, emptyAddr, testTok, testRef, assert.Error},
		{c, "\n", toAddr, emptyAddr, testTok, testRef, assert.Error},
		{c, dne, toAddr, fromAddr, testTok, testRef, assert.Error},
		{c, "\n", toAddr, fromAddr, testTok, testRef, assert.Error},
	}
	for _, test := range tests {
		tmpl := newTemplate("", test)
		assert.NotNil(t, tmpl)
		err := LoadTemplate(c, tmpl)
		test.Err(t, err)
	}
	test := testCase{c, "", emailAddr, emailAddr, testTok, testRef, nil}
	for _, sub := range []string{tmplSub, confSub} {
		tmpl := newTemplate(sub, test)
		require.NotNil(t, tmpl)
		err := LoadTemplate(c, tmpl)
		require.NoError(t, err)
		if sub == "" && tmpl.Subject() != "" {
			assert.NotEqual(t, sub, tmpl.Subject())
			assert.NotEmpty(t, tmpl.Subject())

		} else {
			assert.Equal(t, sub, tmpl.Subject())
		}
	}
}

func TestEmailTemplate_Valid(t *testing.T) {
	t.Parallel()
	e := MailTemplate{}
	assert.Error(t, e.Valid())
	e.Prod.Name = "test"
	assert.Error(t, e.Valid())
	e.to = nameAddr
	assert.Error(t, e.Valid())
	e.to = emailAddr
	assert.NoError(t, e.Valid())
}

func TestEmailTemplate_LoadSender(t *testing.T) {
	t.Parallel()
	require.FileExists(t, testLayout)
	static := config.Mail{
		MailFormat: config.MailFormat{
			Name: confName,
			Link: confLink,
			Logo: testLogo,
		},
	}
	sender := config.Mail{
		MailFormat: config.MailFormat{
			Link: tmplLink,
			Logo: testLogo,
			From: tmplFrom,
		},
		MailTemplates: config.MailTemplates{
			Layout: testLayout,
		},
	}
	tests := []struct {
		conf    config.Mail
		link    string
		service string
		logoURL string
		logo    string
	}{
		{static, confLink, confName, testLogo, path.Base(testLogo)},
		{sender, tmplLink, tmplName, tmplLogo, tmplLogo},
	}
	for _, test := range tests {
		tmpl := &MailTemplate{}
		err := tmpl.LoadLayout(test.conf)
		assert.NoError(t, err)
		assert.Equal(t, test.service, tmpl.Service())
		assert.Equal(t, test.logoURL, tmpl.Logo())
		assert.Equal(t, test.link, tmpl.Prod.Link)
		assert.Equal(t, test.logo, tmpl.Prod.Logo)
		assert.NotEmpty(t, tmpl.Prod.Copyright)
		assert.NotEmpty(t, tmpl.Prod.TroubleText)
	}
	tmpl := &MailTemplate{}
	err := tmpl.LoadLayout(static)
	assert.NoError(t, err)
	assert.Equal(t, static.Link, tmpl.Prod.Link)
	tmpl = &MailTemplate{}
	// bad path
	bad1 := config.Mail{}
	bad1.Layout = "does-not-exist"
	err = tmpl.LoadLayout(bad1)
	assert.Error(t, err)
	// bad yaml
	bad2 := config.Mail{}
	bad2.Layout = testLogo
	err = tmpl.LoadLayout(bad2)
	assert.Error(t, err)
}

func TestEmailTemplate_LoadBody(t *testing.T) {
	t.Parallel()
	require.FileExists(t, testBody)
	static := config.MailTemplate{
		Subject:    confSub,
		LinkFormat: linkFmt,
	}
	tmpl := config.MailTemplate{
		Template:   testBody,
		LinkFormat: linkFmt,
	}
	tests := []struct {
		conf    config.MailTemplate
		to      mail.Address
		subject string
		action  string
		link    string
	}{
		{static, emptyAddr, confSub, testAction, testLink2},
		{static, nameAddr, confSub, "", testTok},
		{tmpl, emailAddr, tmplSub, testAction, testLink2},
		{tmpl, toAddr, tmplSub, "", testTok},
	}
	for _, test := range tests {
		e := MailTemplate{}
		e.Configure(test.conf, test.to, testTok, testRef)
		err := e.LoadBody(test.action, e.Config())
		assert.NoError(t, err)
		assert.Equal(t, test.subject, e.Subject())
		assert.Equal(t, strings.Title(test.to.Name), e.Body.Name)
		assert.Equal(t, test.to, e.to)
		if test.action != "" {
			assert.Equal(t, test.link, e.Body.Actions[0].Button.Link)
		} else {
			assert.Equal(t, test.link, e.Body.Actions[0].InviteCode)
		}
	}
	e := MailTemplate{}
	// bad path
	bad := config.MailTemplate{}
	bad.Template = "does-not-exist"
	err := e.LoadBody(testAction, bad)
	assert.Error(t, err)
	// bad yaml
	bad = config.MailTemplate{}
	bad.Template = testLogo
	err = e.LoadBody(testAction, bad)
	assert.Error(t, err)
	// bad fragment
	bad.LinkFormat = "\n"
	err = e.LoadBody(testAction, bad)
	assert.Error(t, err)
	// bad url
	e = MailTemplate{
		link: "\n",
	}
	err = e.LoadBody(testAction, bad)
	assert.Error(t, err)
}

func testConfig(t *testing.T) *config.Config {
	const tmplFrom = "admin@example.com"
	c := &config.Config{}
	require.FileExists(t, testLayout)
	c.Mail = config.Mail{
		MailFormat: config.MailFormat{
			From: tmplFrom,
		},
		MailTemplates: config.MailTemplates{
			Layout: testLayout,
		},
	}
	require.NotNil(t, c)
	return c
}

func testConfigTemplate(t *testing.T) config.MailTemplate {
	const linkFormat = "/:action/#token=:token"
	// this is actually broken on purpose
	require.FileExists(t, testBody)
	return config.MailTemplate{
		Template:   testBody,
		LinkFormat: linkFormat,
	}
}

func TestEmailTemplate_LoadTemplate(t *testing.T) {
	t.Parallel()
	c := testConfig(t)
	tc := testConfigTemplate(t)
	e := MailTemplate{}
	e.Configure(tc, toAddr, testTok, testRef)
	err := LoadTemplate(c.Mail, &e)
	assert.NoError(t, err)
	e.Configure(tc, toAddr, testTok, "\n")
	err = LoadTemplate(c.Mail, &e)
	assert.Error(t, err)
	c.Mail.Layout = "\n"
	err = LoadTemplate(c.Mail, &e)
	assert.Error(t, err)
	e.Configure(tc, toAddr, testTok, testRef)
	tc.Template = "\n"
	err = LoadTemplate(c.Mail, &e)
	assert.Error(t, err)
}

func TestEmailTemplate_Content(t *testing.T) {
	t.Parallel()
	loadTemplate := func(t *testing.T, tc testCase) MailTemplate {
		tmpl := testConfigTemplate(t)
		e := MailTemplate{}
		e.Configure(tmpl, tc.to, tc.tok, tc.ref)
		return e
	}
	testTemplateContent(t, func(tc testCase) (string, Template) {
		e := loadTemplate(t, tc)
		return "email_template_file", &e
	})
}

func testTemplateContent(t *testing.T, newTemplate func(test testCase) (string, Template)) {
	c := tconf.Config(t).Mail
	c.Logo = testLogo
	flat := c
	flat.Theme = "flat"
	tests := []testCase{
		{c, "", toAddr, emptyAddr, testTok, testRef, assert.NoError},
		{c, "", toAddr, fromAddr, testTok, testRef, assert.NoError},
		{flat, "", toAddr, emptyAddr, testTok, testRef, assert.NoError},
		{flat, "", toAddr, fromAddr, testTok, testRef, assert.NoError},
	}
	for _, test := range tests {
		action, tmpl := newTemplate(test)
		assert.NotNil(t, tmpl)
		err := LoadTemplate(test.mc, tmpl)
		assert.NoError(t, err)
		html, err := tmpl.HTML()
		assert.NoError(t, err)
		action = strings.ReplaceAll(action, "/", "_")
		htmlFile := fmt.Sprintf("%s%s_%s_test", testDir, action, test.mc.Theme)
		assertEqualContent(t, htmlFile+".html", html)
		plain, err := tmpl.PlainText()
		assert.NoError(t, err)
		plainFile := fmt.Sprintf("%s%s_test", testDir, action)
		assertEqualContent(t, plainFile+".txt", plain)
	}
}

func TestEmail_Load(t *testing.T) {
	t.Parallel()
	mc := tconf.Config(t).Mail
	mc.Logo = testLogo
	body, err := loadTemplateFile(testBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	c := config.MailTemplate{}
	c.Subject = confSub
	em := NewEmail(c, toAddr, body)
	require.NotNil(t, em)
	err = LoadTemplate(mc, em)
	assert.NoError(t, err)
}

func TestEmail_Content(t *testing.T) {
	t.Parallel()
	body, err := loadTemplateFile(testBody)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	testTemplateContent(t, func(tc testCase) (string, Template) {
		c := config.MailTemplate{}
		c.Subject = confSub
		e := NewEmail(c, tc.to, body)
		return "email_template_body", e
	})
}
