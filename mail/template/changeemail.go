package template

import (
	"fmt"
	"net/mail"

	"github.com/jrapoport/gothic/config"
	"github.com/matcornic/hermes/v2"
)

// ChangeEmailAction confirm email action
const ChangeEmailAction = "change/email"

// ChangeEmail mail template
type ChangeEmail struct {
	MailTemplate
	newAddress string
}

var _ Template = (*ChangeEmail)(nil)

// NewChangeEmail returns a new change mail template.
func NewChangeEmail(c config.MailTemplate, to mail.Address, newAddress, token, referralURL string) *ChangeEmail {
	e := new(ChangeEmail)
	e.newAddress = newAddress
	e.Configure(c, to, token, referralURL)
	return e
}

// Action returns the action for the mail template.
func (e ChangeEmail) Action() string {
	return ChangeEmailAction
}

// Subject returns the subject for the mail.
func (e ChangeEmail) Subject() string {
	if e.MailTemplate.Subject() != "" {
		return e.MailTemplate.Subject()
	}
	return e.subject()
}

// LoadBody loads the body for the mail.
func (e *ChangeEmail) LoadBody(action string, tc config.MailTemplate) error {
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
	e.Body.Outros = append([]string{e.outro()}, e.Body.Outros...)
	return nil
}

func (e ChangeEmail) subject() string {
	const subjectFormat = "%s email change request"
	return fmt.Sprintf(subjectFormat, e.Service())
}

func (e ChangeEmail) intro() string {
	const introFormat = "You received this message because there was a request " +
		"to change the email address you use to access your %s account."
	return fmt.Sprintf(introFormat, e.Service())
}

func (e ChangeEmail) instructions() string {
	const instructFormat = "Please click the button below to confirm and change" +
		" your email to: %s"
	return fmt.Sprintf(instructFormat, e.newAddress)
}

func (e ChangeEmail) buttonText() string {
	return "Change Email"
}

func (e ChangeEmail) outro() string {
	return "Once confirmed, your login email will change to the new address. If" +
		" you did not request this change, no further action is required. You " +
		"can safely ignore this message."
}
