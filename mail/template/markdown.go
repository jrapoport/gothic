package template

import (
	"net/mail"

	"github.com/jrapoport/gothic/config"
	"github.com/matcornic/hermes/v2"
)

// Markdown is a wrapper for a generic markdown mail template
type Markdown struct {
	MailTemplate
}

var _ Template = (*Markdown)(nil)

// NewMarkdownMail returns a new markdown email
func NewMarkdownMail(c config.MailTemplate, to mail.Address, markdown string) *Markdown {
	md := new(Markdown)
	md.Body.FreeMarkdown = hermes.Markdown(markdown)
	md.Configure(c, to, "", "")
	return md
}

// Action returns the action for the mail template.
func (md Markdown) Action() string {
	return "markdown"
}

// LoadBody loads the body for the mail.
func (md *Markdown) LoadBody(_ string, tc config.MailTemplate) error {
	err := md.MailTemplate.LoadBody("", tc)
	if err != nil {
		return err
	}
	return nil
}
