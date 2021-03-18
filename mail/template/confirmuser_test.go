package template

import "testing"

func TestConfirmUser_Load(t *testing.T) {
	t.Parallel()
	testTemplateLoad(t, func(sub string, test testCase) Template {
		c := test.mc.ConfirmUser
		c.Subject = sub
		c.Template = test.tmpl
		return NewConfirmUser(c, test.to, test.tok, test.ref)
	})
}

func TestConfirmUser_Content(t *testing.T) {
	t.Parallel()
	testTemplateContent(t, func(test testCase) (string, Template) {
		e := NewConfirmUser(test.mc.ConfirmUser, test.to, test.tok, test.ref)
		return e.Action(), e
	})
}
