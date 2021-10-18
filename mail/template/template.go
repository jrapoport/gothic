package template

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/mail"
	"path/filepath"
	"strings"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/validate"
	"github.com/jrapoport/gothic/utils"
	"github.com/matcornic/hermes/v2"
	"gopkg.in/yaml.v3"
)

// Template is the interface for mail templates.
type Template interface {
	Action() string
	To() string
	Logo() string
	Subject() string
	Valid() error
	HTML() (string, error)
	PlainText() (string, error)
	Config() config.MailTemplate
	LoadLayout(c config.Mail) error
	LoadBody(action string, c config.MailTemplate) error
}

// LoadTemplate loads a mail template.
func LoadTemplate(c config.Mail, t Template) error {
	err := t.LoadLayout(c)
	if err != nil {
		return err
	}
	err = t.LoadBody(t.Action(), t.Config())
	if err != nil {
		return err
	}
	return t.Valid()
}

// MailTemplate holds a mail template.
type MailTemplate struct {
	theme        hermes.Theme
	Prod         hermes.Product
	Body         hermes.Body
	config       config.MailTemplate
	to           mail.Address
	subject      string
	logo         string
	link         string
	token        string
	bodyTemplate string
}

var _ Template = (*MailTemplate)(nil)

// Configure configures a mail template.
func (e *MailTemplate) Configure(c config.MailTemplate, to mail.Address, token, linkURL string) {
	e.config = c
	e.SetToAddress(to)
	e.token = token
	e.link = linkURL
}

// NewEmail returns a new configured email template
func NewEmail(c config.MailTemplate, to mail.Address, body string) *MailTemplate {
	e := new(MailTemplate)
	e.bodyTemplate = body
	e.Configure(c, to, "", "")
	return e
}

// Config returns the config used t init a mail template.
func (e MailTemplate) Config() config.MailTemplate {
	return e.config
}

// LoadLayout loads the mail template layout.
func (e *MailTemplate) LoadLayout(c config.Mail) error {
	switch c.Theme {
	case "flat":
		e.theme = new(hermes.Flat)
	default:
		e.theme = new(hermes.Default)
	}
	if e.link == "" {
		// only set this IF the referral URL a.k.a.
		// the current value of link, is empty.
		e.link = c.Link
	}
	// sets e.Prod.Name
	e.SetService(c.Name)
	e.Prod.Link = c.Link
	e.Prod.Logo = c.Logo
	e.Prod.Copyright = e.defaultCopyright()
	e.Prod.TroubleText = e.defaultHelp()
	tmpl, err := loadTemplateFile(c.Layout)
	if err != nil {
		err = fmt.Errorf("layout file: %w", err)
		return err
	}
	if tmpl != "" {
		err = yaml.Unmarshal([]byte(tmpl), &e.Prod)
		if err != nil {
			err = fmt.Errorf("unmarshal layout template: %w", err)
			return err
		}
	}
	e.SetLogo(e.Prod.Logo)
	return err
}

// LoadBody loads the mail template body.
func (e *MailTemplate) LoadBody(action string, ct config.MailTemplate) (err error) {
	e.link, err = e.generateLink(action, ct)
	if err != nil {
		return err
	}
	e.subject = ct.Subject
	e.Body.Signature = e.defaultSignature()
	e.Body.Outros = []string{e.defaultOutro()}
	if e.token != "" {
		a := hermes.Action{}
		if action != "" {
			a.Button.Link = e.link
		} else {
			a.InviteCode = e.token
		}
		e.Body.Actions = []hermes.Action{a}
	}
	if e.bodyTemplate == "" {
		e.bodyTemplate, err = loadTemplateFile(ct.Template)
		if err != nil {
			err = fmt.Errorf("body file: %w", err)
			return err
		}
	}
	if e.bodyTemplate != "" {
		err = yaml.Unmarshal([]byte(e.bodyTemplate), &e.Body)
		if err != nil {
			err = fmt.Errorf("unmarshal body template: %w", err)
			return err
		}
	}
	if len(e.Body.Actions) <= 0 {
		return nil
	}
	a := &e.Body.Actions[0]
	if e.token == "" {
		return nil
	}
	if action != "" {
		a.Button.Link = e.link
	} else {
		a.InviteCode = e.token
	}
	return nil
}

// Action returns the action for the mail template (e.g. confirm).
func (e MailTemplate) Action() string {
	return ""
}

// Token returns the token associated with the mail template (if any)
func (e MailTemplate) Token() string {
	return e.token
}

// SetToAddress sets the to: address of the email.
func (e *MailTemplate) SetToAddress(address mail.Address) {
	e.to = address
	e.Body.Name = strings.Title(address.Name)
}

// To returns the to: address of the email.
func (e MailTemplate) To() string {
	return e.to.String()
}

// Subject returns the subject address of the email.
func (e MailTemplate) Subject() string {
	return e.subject
}

// SetLogo sets the logo url or path for the mail header.
func (e *MailTemplate) SetLogo(logo string) {
	if utils.IsLocalPath(logo) {
		e.Prod.Logo = filepath.Base(logo)
	}
	e.logo = logo
}

// Logo returns the logo url or path for the mail header.
func (e MailTemplate) Logo() string {
	return e.logo
}

// Link returns the service url for the mail header.
func (e MailTemplate) Link() string {
	return e.link
}

// SetService sets the name of the service for the mail header.
func (e *MailTemplate) SetService(n string) {
	e.Prod.Name = n
}

// Service returns the name of the service for the mail header.
func (e MailTemplate) Service() string {
	return e.Prod.Name
}

// Valid returns nil if the template is valid.
func (e MailTemplate) Valid() error {
	if e.Service() == "" {
		return errors.New("service invalid")
	}
	if _, err := validate.Email(e.to.Address); err != nil {
		err = fmt.Errorf("to email %w", err)
		return err
	}
	return nil
}

// HTML returns an html version of the email template.
func (e MailTemplate) HTML() (string, error) {
	h := e.generator()
	return h.GenerateHTML(e.email())
}

// PlainText returns a plaintext version of the email template.
func (e MailTemplate) PlainText() (string, error) {
	h := e.generator()
	return h.GeneratePlainText(e.email())
}

func (e MailTemplate) generator() hermes.Hermes {
	return hermes.Hermes{
		Theme:   e.theme,
		Product: e.Prod,
	}
}

func (e MailTemplate) email() hermes.Email {
	return hermes.Email{Body: e.Body}
}

func (e MailTemplate) generateLink(action string, ct config.MailTemplate) (string, error) {
	link, err := utils.JoinLink(e.link, ct.LinkFormat)
	if err != nil {
		return "", err
	}
	link = config.FormatLink(link, action, e.token)
	return utils.NormalizeURL(link)
}

func (e MailTemplate) defaultCopyright() string {
	const copyright = "Copyright © %d %s"
	return fmt.Sprintf(copyright, time.Now().UTC().Year(), e.Service())
}

func (e MailTemplate) defaultSignature() string {
	return "Thanks"
}

func (e MailTemplate) defaultOutro() string {
	return "Need help, or have questions? Please contact support. Do not reply " +
		"to this email."
}

func (e MailTemplate) defaultHelp() string {
	return "If the \"{ACTION}\" button is not working for you, just copy and " +
		"paste the URL below into your web browser."
}

func loadTemplateFile(file string) (string, error) {
	if file == "" {
		return "", nil
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
