package template

import (
	"fmt"
	"net/mail"

	"github.com/jrapoport/gothic/config"
	"github.com/matcornic/hermes/v2"
)

// InviteUserAction invite user action
const InviteUserAction = "invite"

// InviteUser mail template
type InviteUser struct {
	MailTemplate
	fromName string
}

var _ Template = (*InviteUser)(nil)

// NewInviteUser returns a new user invite mail template
func NewInviteUser(c config.MailTemplate, from, to mail.Address, token, referralURL string) *InviteUser {
	e := new(InviteUser)
	e.fromName = from.Name
	e.Configure(c, to, token, referralURL)
	return e
}

// Action returns the action for the mail template.
func (e InviteUser) Action() string {
	return InviteUserAction
}

// Subject returns the subject for the mail.
func (e InviteUser) Subject() string {
	if e.MailTemplate.Subject() != "" {
		return e.MailTemplate.Subject()
	}
	return e.subject()
}

// LoadBody loads the body for the mail.
func (e *InviteUser) LoadBody(action string, tc config.MailTemplate) error {
	err := e.MailTemplate.LoadBody(action, tc)
	if err != nil {
		return err
	}
	if len(e.Body.Intros) <= 0 {
		e.Body.Intros = e.intros()
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

func (e InviteUser) subject() string {
	if e.fromName == "" {
		const anonFormat = "Your invite to %s"
		return fmt.Sprintf(anonFormat, e.Service())
	}
	const subjectFormat = "%s invited you to %s"
	return fmt.Sprintf(subjectFormat, e.fromName, e.Service())
}

func (e InviteUser) intros() []string {
	const intro1Format = "You have been invited to sign up for %s!"
	const intro2Format = "You received this message because you have been " +
		"invited to join your friends on %s."
	return []string{
		fmt.Sprintf(intro1Format, e.Service()),
		fmt.Sprintf(intro2Format, e.Prod.Link),
	}
}

func (e InviteUser) instructions() string {
	return "To accept your invite and create your account, " +
		"please click the button below:"
}

func (e InviteUser) buttonText() string {
	return "Accept Invite"
}
