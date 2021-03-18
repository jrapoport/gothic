package template

import "testing"

func TestSignupCode_Load(t *testing.T) {
	t.Parallel()
	testTemplateLoad(t, func(sub string, test testCase) Template {
		c := test.mc.SignupCode
		c.Subject = sub
		c.Template = test.tmpl
		return NewSignupCode(c, test.from, test.to, test.tok, test.ref)
	})
}

func TestSignupCode_Content(t *testing.T) {
	t.Parallel()
	const testCode = "123456"
	testTemplateContent(t, func(tc testCase) (string, Template) {
		e := NewSignupCode(tc.mc.SignupCode, tc.from, tc.to, testCode, tc.ref)
		return e.Action(), e
	})
}
