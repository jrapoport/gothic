package template

import "testing"

func TestInvite_Load(t *testing.T) {
	t.Parallel()
	testTemplateLoad(t, func(sub string, test testCase) Template {
		c := test.mc.InviteUser
		c.Subject = sub
		c.Template = test.tmpl
		return NewInviteUser(c, test.from, test.to, test.tok, test.ref)
	})
}

func TestInvite_Content(t *testing.T) {
	t.Parallel()
	testTemplateContent(t, func(test testCase) (string, Template) {
		e := NewInviteUser(test.mc.InviteUser, test.from, test.to, test.tok, test.ref)
		return e.Action(), e
	})
}
