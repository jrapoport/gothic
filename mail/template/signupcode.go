package template

import (
	"fmt"
	"net/mail"

	"github.com/jrapoport/gothic/config"
	"github.com/matcornic/hermes/v2"
)

// SignupCodeAction is the signup code action.
const SignupCodeAction = "signup"

// SignupCode mail template
type SignupCode struct {
	InviteUser
}

var _ Template = (*SignupCode)(nil)

// NewSignupCode returns a new signup code mail template
func NewSignupCode(c config.MailTemplate, from, to mail.Address, token, referralURL string) *SignupCode {
	e := new(SignupCode)
	e.InviteUser = *NewInviteUser(c, from, to, token, referralURL)
	return e
}

// Action returns the action for the mail template.
func (e SignupCode) Action() string {
	return SignupCodeAction
}

// LoadBody loads the body for the mail.
func (e *SignupCode) LoadBody(action string, tc config.MailTemplate) error {
	err := e.InviteUser.LoadBody(action, tc)
	if err != nil {
		return err
	}
	if len(e.Body.Actions) <= 0 {
		e.Body.Actions = append(e.Body.Actions, hermes.Action{})
	}
	a := &e.Body.Actions[0]
	a.Button = hermes.Button{}
	a.Instructions = e.instructions()
	a.InviteCode = e.Token()
	return nil
}

func (e SignupCode) instructions() string {
	const instructFormat = "To get started with %s," +
		" please visit our signup page and enter the code below: %s"
	return fmt.Sprintf(instructFormat, e.Service(), e.Link())
}
