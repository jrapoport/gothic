package template

import (
	"testing"

	"github.com/jrapoport/gothic/config"
)

const markdownText = `
Title: Newsletter Number 6
Date: 12-9-2019 10:04am
Template: newsletter
URL: newsletter/issue-6.html
save_as: newsletter/issue-6.html

Welcome to the 6th edition of this newsletter.

## Around the site
Hello Subscriber
`

func TestMarkdownEmail_Load(t *testing.T) {
	t.Parallel()
	testTemplateLoad(t, func(sub string, test testCase) Template {
		c := config.MailTemplate{}
		c.Subject = sub
		c.Template = test.tmpl
		return NewMarkdownMail(c, test.to, markdownText)
	})
}

func TestMarkdownEmail_Content(t *testing.T) {
	t.Parallel()
	testTemplateContent(t, func(tc testCase) (string, Template) {
		e := NewMarkdownMail(config.MailTemplate{
			Subject: "Markdown Subject",
		}, tc.to, markdownText)
		return e.Action(), e
	})
}
