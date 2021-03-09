package template

import "testing"

func TestResetPassword_Load(t *testing.T) {
	testTemplateLoad(t, func(sub string, test testCase) Template {
		c := test.mc.ResetPassword
		c.Subject = sub
		c.Template = test.tmpl
		return NewResetPassword(c, test.to, test.tok, test.ref)
	})
}

func TestResetPassword_Content(t *testing.T) {
	testTemplateContent(t, func(tc testCase) (string, Template) {
		e := NewResetPassword(tc.mc.ResetPassword, tc.to, tc.tok, tc.ref)
		return e.Action(), e
	})
}
