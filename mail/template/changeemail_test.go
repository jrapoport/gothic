package template

import "testing"

func TestChangeEmail_Load(t *testing.T) {
	testTemplateLoad(t, func(sub string, test testCase) Template {
		c := test.mc.ChangeEmail
		c.Subject = sub
		c.Template = test.tmpl
		return NewChangeEmail(c, test.to, fromEmail, test.tok, test.ref)
	})
}

func TestChangeEmail_Content(t *testing.T) {
	testTemplateContent(t, func(test testCase) (string, Template) {
		e := NewChangeEmail(test.mc.ChangeEmail, test.to, fromEmail, test.tok, test.ref)
		return e.Action(), e
	})
}
