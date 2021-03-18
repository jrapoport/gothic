package core

import (
	"testing"
	"time"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/codes"
	"github.com/jrapoport/gothic/core/tokens"
	"github.com/jrapoport/gothic/core/tokens/jwt"
	"github.com/jrapoport/gothic/mail/template"
	"github.com/jrapoport/gothic/models/code"
	"github.com/jrapoport/gothic/models/token"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store/types/key"
	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/jrapoport/gothic/test/tconf"
	"github.com/jrapoport/gothic/test/tutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mailAPI(t *testing.T) *API {
	var c = tconf.TempDB(t)
	c.Mail.Host = ""
	c.Validation.PasswordRegex = ""
	c.Signup.AutoConfirm = false
	a, err := NewAPI(c)
	require.NoError(t, err)
	require.NotNil(t, a)
	t.Cleanup(func() {
		err = a.Shutdown()
		require.NoError(t, err)
	})
	return a
}

func forceExtProvider(t *testing.T, a *API, u *user.User) {
	u.Provider = provider.Google
	err := a.conn.Save(u).Error
	require.NoError(t, err)
}

func TestAPI_OpenMail(t *testing.T) {
	a := mailAPI(t)
	assert.True(t, a.mail.IsOffline())
	err := a.OpenMail()
	assert.NoError(t, err)
	assert.True(t, a.mail.IsOffline())
	a.config, _ = tconf.MockSMTP(t, a.config)
	err = a.OpenMail()
	assert.NoError(t, err)
	assert.False(t, a.mail.IsOffline())
	a.CloseMail()
	assert.Nil(t, a.mail)
	err = a.OpenMail()
	assert.NoError(t, err)
	a.CloseMail()
	assert.Nil(t, a.mail)
}

func TestAPI_ValidateEmail(t *testing.T) {
	a := mailAPI(t)
	assert.True(t, a.mail.IsOffline())
	e := tutils.RandomEmail()
	_, err := a.ValidateEmail(e)
	assert.NoError(t, err)
	_, err = a.ValidateEmail("")
	assert.Error(t, err)
	a.config, _ = tconf.MockSMTP(t, a.config)
	err = a.OpenMail()
	assert.NoError(t, err)
	_, err = a.ValidateEmail(e)
	// this is only an error bc it is a mock smtp server
	assert.Error(t, err)
	_, err = a.ValidateEmail("")
	assert.Error(t, err)
	a.CloseMail()
	assert.Nil(t, a.mail)
	_, err = a.ValidateEmail(e)
	assert.NoError(t, err)
	_, err = a.ValidateEmail("")
	assert.Error(t, err)
}

func TestAPI_SendConfirmUser(t *testing.T) {
	const referral = "http://www.example.com/createUser/user"
	a := mailAPI(t)
	ctx := testContext(a)
	a.config.Mail.ConfirmUser.ReferralURL = referral
	u := testUser(t, a)
	assert.False(t, u.IsConfirmed())
	testSend(t, a, u, template.ConfirmUserAction,
		func() error {
			return a.SendConfirmUser(ctx, u.ID)
		},
		func(tok string) {
			ct, err := tokens.GetConfirmToken(a.conn, tok)
			assert.NoError(t, err)
			assert.Equal(t, u.ID, ct.UserID)
		},
		func() {
			// rate limit
			var ct token.ConfirmToken
			err := a.conn.First(&ct, "user_id = ?", u.ID).Error
			assert.NoError(t, err)
			now := time.Now().UTC()
			ct.SentAt = &now
			err = a.conn.Save(ct).Error
			assert.NoError(t, err)
		})
	confirmUser(t, a, u)
	err := a.SendConfirmUser(ctx, u.ID)
	assert.NoError(t, err)
	forceExtProvider(t, a, u)
	err = a.SendConfirmUser(ctx, u.ID)
	assert.Error(t, err)
	banUser(t, a, u)
	err = a.SendConfirmUser(ctx, u.ID)
	assert.Error(t, err)
	a.mail = nil
	err = a.SendConfirmUser(ctx, u.ID)
	assert.NoError(t, err)
}

func TestAPI_SendResetPassword(t *testing.T) {
	const referral = "http://www.example.com/password/"
	a := mailAPI(t)
	ctx := testContext(a)
	a.config.Mail.ResetPassword.ReferralURL = referral
	u := testUser(t, a)
	assert.False(t, u.IsConfirmed())
	act := template.ResetPasswordAction
	testSend(t, a, u, act,
		func() error {
			return a.SendResetPassword(ctx, u.ID)
		},
		func(tok string) {
			ct, err := tokens.GetConfirmToken(a.conn, tok)
			assert.NoError(t, err)
			require.NotNil(t, ct)
			assert.Equal(t, u.ID, ct.UserID)
		},
		func() {
			// rate limit
			var ct token.ConfirmToken
			err := a.conn.First(&ct, "user_id = ?", u.ID).Error
			assert.NoError(t, err)
			now := time.Now().UTC()
			ct.SentAt = &now
			err = a.conn.Save(ct).Error
			assert.NoError(t, err)
		})
	forceExtProvider(t, a, u)
	err := a.SendResetPassword(ctx, u.ID)
	assert.Error(t, err)
	banUser(t, a, u)
	err = a.SendResetPassword(ctx, u.ID)
	assert.Error(t, err)
	a.mail = nil
	err = a.SendResetPassword(ctx, u.ID)
	assert.NoError(t, err)
}

func TestAPI_SendChangeEmail(t *testing.T) {
	const referral = "http://www.example.com/change/"
	a := mailAPI(t)
	ctx := testContext(a)
	a.config.Mail.ChangeEmail.ReferralURL = referral
	to := tutils.RandomEmail()
	u := testUser(t, a)
	u = confirmUser(t, a, u)
	act := template.ChangeEmailAction
	testSend(t, a, u, act,
		func() error {
			return a.SendChangeEmail(ctx, u.ID, to)
		},
		func(tok string) {
			data, err := jwt.ParseData(a.config.JWT, tok)
			assert.NoError(t, err)
			assert.Equal(t, to, data[key.Email])
			tok, ok := data[key.Token].(string)
			assert.True(t, ok)
			ct, err := tokens.GetConfirmToken(a.conn, tok)
			assert.NoError(t, err)
			assert.Equal(t, u.ID, ct.UserID)
		},
		func() {
			// rate limit
			var ct token.ConfirmToken
			err := a.conn.First(&ct, "user_id = ?", u.ID).Error
			assert.NoError(t, err)
			now := time.Now().UTC()
			ct.SentAt = &now
			err = a.conn.Save(ct).Error
			assert.NoError(t, err)
		})
	u = testUser(t, a)
	assert.False(t, u.IsConfirmed())
	err := a.SendChangeEmail(ctx, u.ID, to)
	assert.Error(t, err)
	forceExtProvider(t, a, u)
	err = a.SendChangeEmail(ctx, u.ID, to)
	assert.Error(t, err)
	banUser(t, a, u)
	err = a.SendChangeEmail(ctx, u.ID, to)
	assert.Error(t, err)
	a.mail = nil
	err = a.SendChangeEmail(ctx, u.ID, to)
	assert.NoError(t, err)
}

func TestAPI_SendInviteUser(t *testing.T) {
	a := mailAPI(t)
	ctx := testContext(a)
	a.config.Signup.Invites = config.Users
	u := testUser(t, a)
	to := tutils.RandomEmail()
	confirmUser(t, a, u)
	act := template.InviteUserAction
	testSend(t, a, u, act,
		func() error {
			return a.SendInviteUser(ctx, u.ID, to)
		},
		func(tok string) {
			data, err := jwt.ParseData(a.config.JWT, tok)
			assert.NoError(t, err)
			assert.Equal(t, to, data[key.Email])
			tok, ok := data[key.Token].(string)
			assert.True(t, ok)
			sc, err := codes.GetUsableSignupCode(a.conn, tok)
			assert.NoError(t, err)
			assert.Equal(t, u.ID, sc.UserID)
		},
		func() {
			// rate limit
			var sc code.SignupCode
			err := a.conn.First(&sc, "user_id = ?", u.ID).Error
			assert.NoError(t, err)
			now := time.Now().UTC()
			sc.SentAt = &now
			err = a.conn.Save(sc).Error
			assert.NoError(t, err)
		})

	// no admin rate limit
	adm := promoteUser(t, a, u)
	err := a.SendInviteUser(ctx, adm.ID, to)
	assert.NoError(t, err)
	// admin restriction
	a.config.Signup.Invites = config.Admins
	err = a.SendInviteUser(ctx, adm.ID, to)
	assert.NoError(t, err)
	// bad mail
	err = a.SendInviteUser(ctx, u.ID, "@")
	assert.Error(t, err)
	// not an admin
	u = testUser(t, a)
	// not confirmed
	err = a.SendInviteUser(ctx, u.ID, tutils.RandomEmail())
	assert.Error(t, err)
	confirmUser(t, a, u)
	err = a.SendInviteUser(ctx, u.ID, to)
	assert.Error(t, err)
	// banned
	banUser(t, a, adm)
	err = a.SendInviteUser(ctx, adm.ID, to)
	assert.Error(t, err)
	// signup disabled
	a.config.Signup.Disabled = true
	err = a.SendInviteUser(ctx, adm.ID, to)
	assert.Error(t, err)
	// no limits
	err = a.SendInviteUser(ctx, user.SystemID, to)
	assert.NoError(t, err)
}

func testSend(t *testing.T, a *API, u *user.User, action string,
	send func() error,
	testToken func(tok string),
	rateLimit func(),
) {
	t.Cleanup(func() {
		a.CloseMail()
		err := send()
		assert.NoError(t, err)
	})
	// offline
	err := send()
	assert.NoError(t, err)
	// online validation fail
	var mock *tconf.SMTPMock
	a.config, mock = tconf.MockSMTP(t, a.config)
	a.config.Mail.SendLimit = 0
	a.config.Mail.SpamProtection = true
	err = a.OpenMail()
	assert.NoError(t, err)
	err = send()
	assert.Error(t, err)
	a.config.Mail.SpamProtection = false
	err = a.OpenMail()
	assert.NoError(t, err)
	// online
	var tok string
	mock.AddHook(t, func(email string) {
		tok = tconf.GetEmailToken(action, email)
	})
	err = send()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return tok != ""
	}, 1*time.Second, 10*time.Millisecond)
	testToken(tok)
	// rate limit
	if rateLimit != nil {
		rateLimit()
		a.config.Mail.SendLimit = 5 * time.Minute
		err = send()
		if u.Role == user.RoleUser {
			assert.ErrorIs(t, err, config.ErrRateLimitExceeded)
		} else {
			assert.NoError(t, err)
		}
	}
}
