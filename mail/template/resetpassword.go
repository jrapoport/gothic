package template

import (
	"fmt"
	"net/mail"

	"github.com/jrapoport/gothic/config"
	"github.com/matcornic/hermes/v2"
)

// ResetPasswordAction reset password action
const ResetPasswordAction = "reset/password"

// ResetPassword mail template
type ResetPassword struct {
	MailTemplate
}

var _ Template = (*ResetPassword)(nil)

// NewResetPassword returns a new reset password email
func NewResetPassword(c config.MailTemplate, to mail.Address, token, referralURL string) *ResetPassword {
	e := new(ResetPassword)
	e.Configure(c, to, token, referralURL)
	return e
}

// Action returns the action for the mail template.
func (e ResetPassword) Action() string {
	return ResetPasswordAction
}

// Subject returns the subject for the mail.
func (e ResetPassword) Subject() string {
	if e.MailTemplate.Subject() != "" {
		return e.MailTemplate.Subject()
	}
	return e.subject()
}

// LoadBody loads the body for the mail.
func (e *ResetPassword) LoadBody(action string, tc config.MailTemplate) error {
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

func (e ResetPassword) subject() string {
	return "Reset your password"
}

func (e ResetPassword) intro() string {
	const introFormat = "You received this message because there was a request " +
		"to reset the password for your %s account."
	return fmt.Sprintf(introFormat, e.Service())
}

func (e ResetPassword) instructions() string {
	return "To reset your password, please click the button below:"
}

func (e ResetPassword) buttonText() string {
	return "Reset Password"
}

func (e ResetPassword) outro() string {
	return "If you did not request this change, no further action is required. " +
		"You can safely ignore this message."
}
