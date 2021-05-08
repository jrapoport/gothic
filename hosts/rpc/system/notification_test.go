package system

import (
	"context"
	"github.com/jrapoport/gothic/mail"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/test/tcore"
	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSubject  = "Test Subject"
	testHTML     = "<html>test_html_notification</html>"
	testMarkdown = `
Title: Newsletter Number 6
Date: 12-9-2019 10:04am
Template: newsletter
URL: newsletter/issue-6.html
save_as: newsletter/issue-6.html

Welcome to the 6th edition of this newsletter.

## Around the site
Hello Markdown
`
	testTemplate = `
intros:
  - Welcome to Gothic!
actions:
  - instructions: 'To a friend know, please click here:'
    button:
      text: Hello Template
      link: https://www.example.com/friend?token=invite-a-friend
outros:
  - We're very excited to have you on board.

`
)

func TestSystemServer_SendEmailNotification(t *testing.T) {
	var testPlain = "test_plain_notification"
	t.Parallel()
	srv := testServer(t)
	ctx := context.Background()
	// no id or email
	req := &system.EmailRequest{}
	_, err := srv.SendEmailNotification(ctx, req)
	assert.Error(t, err)
	// bad id
	req = &system.EmailRequest{UserId: "1"}
	_, err = srv.SendEmailNotification(ctx, req)
	assert.Error(t, err)
	// offline
	req = &system.EmailRequest{
		UserId: uuid.New().String(),
	}
	res, err := srv.SendEmailNotification(ctx, req)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.False(t, res.Sent)
	// mail online
	s, mock := tsrv.RPCServer(t, true)
	s.Config().Signup.AutoConfirm = true
	s.Config().Mail.SMTP.SpamProtection = false
	srv = newSystemServer(s)
	// id not found
	req = &system.EmailRequest{
		UserId: uuid.New().String(),
	}
	_, err = srv.SendEmailNotification(ctx, req)
	assert.Error(t, err)
	// create a user
	u, _ := tcore.TestUser(t, srv.API, "", false)
	require.NotNil(t, u)
	req = &system.EmailRequest{
		UserId: u.ID.String(),
	}
	_, err = srv.SendEmailNotification(ctx, req)
	assert.Error(t, err)
	// success
	tests := map[string]struct {
		typ      mail.Type
		content  string
		expected string
	}{
		"html":     {mail.HTML, testHTML, testHTML},
		"markdown": {mail.Markdown, testMarkdown, "Hello Markdown"},
		"template": {mail.Template, testTemplate, "Hello Template"},
	}
	var recv string
	var mu sync.Mutex
	mock.AddHook(t, func(email string) {
		mu.Lock()
		defer mu.Unlock()
		recv = email
	})
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req = &system.EmailRequest{
				UserId:    u.ID.String(),
				Subject:   testSubject,
				Plaintext: &testPlain,
			}
			switch test.typ {
			case mail.HTML:
				req.Content = &system.EmailRequest_Html{Html: test.content}
			case mail.Markdown:
				req.Content = &system.EmailRequest_Markdown{Markdown: test.content}
			case mail.Template:
				req.Content = &system.EmailRequest_Body{Body: test.content}
			}
			res, err = srv.SendEmailNotification(ctx, req)
			assert.NoError(t, err)
			require.NotNil(t, res)
			assert.True(t, res.Sent)
			assert.Eventually(t, func() bool {
				if !strings.Contains(recv, testSubject) {
					return false
				}
				if !strings.Contains(recv, test.expected) {
					return false
				}
				return true
			}, 1*time.Second, 10*time.Millisecond)
			if test.typ == mail.HTML {
				assert.Contains(t, recv, testPlain)
			}
		})
	}
}
