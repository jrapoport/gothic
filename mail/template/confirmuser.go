package template

import (
	"fmt"
	"net/mail"

	"github.com/jrapoport/gothic/config"
	"github.com/matcornic/hermes/v2"
)

// ConfirmUserAction confirm user action
const ConfirmUserAction = "confirm"

// ConfirmUser mail template
type ConfirmUser struct {
	MailTemplate
}

var _ Template = (*ConfirmUser)(nil)

// NewConfirmUser returns a new user confirmation mail template.
func NewConfirmUser(c config.MailTemplate, to mail.Address, token, referralURL string) *ConfirmUser {
	e := new(ConfirmUser)
	e.Configure(c, to, token, referralURL)
	return e
}

// Action returns the action for the mail template.
func (e ConfirmUser) Action() string {
	return ConfirmUserAction
}

// Subject returns the subject for the mail.
func (e ConfirmUser) Subject() string {
	if e.MailTemplate.Subject() != "" {
		return e.MailTemplate.Subject()
	}
	return e.subject()
}

// LoadBody loads the body for the mail.
func (e *ConfirmUser) LoadBody(action string, tc config.MailTemplate) error {
	err := e.MailTemplate.LoadBody(action, tc)
	if err != nil {
		return err
	}
	if len(e.Body.Intros) <= 0 {
		e.Body.Intros = []string{e.intro()}
	}
	if len(e.Body.Actions) <= 0 {
		e.Body.Actions = append(e.Body.Actions, hermes.Action{})
	}
	a := &e.Body.Actions[0]
	if a.Instructions == "" {
		a.Instructions = e.instructions()
	}
	if a.Button.Text == "" {
		a.Button.Text = e.buttonText()
	}
	if a.Button.Link == "" {
		a.Button.Link = e.Link()
	}
	return nil
}

func (e ConfirmUser) subject() string {
	const subjectFormat = "Confirm your %s account"
	return fmt.Sprintf(subjectFormat, e.Service())
}

func (e ConfirmUser) intro() string {
	const introFormat = "Welcome to %s! We are excited to have you on board."
	return fmt.Sprintf(introFormat, e.Service())
}

func (e ConfirmUser) instructions() string {
	return "Please click the button below to confirm your email address and " +
		"finish setting up your account:"
}

func (e ConfirmUser) buttonText() string {
	return "Confirm Email"
}
