package config

import (
	"net/url"
	"strings"
)

// MailTemplates config
type MailTemplates struct {
	// Theme is theme to use when formatting emails.
	// Options are "default", "flat", "" = default
	Theme string `json:"theme"`
	// MailTemplate is a path to a configuration file with additional customization.
	// For an example, SEE: mail/email/testdata/template_layout.yaml
	Layout string `json:"layout"`

	// ChangeEmail email customization
	ChangeEmail MailTemplate `json:"change_email" yaml:"change_email" mapstructure:"change_email"`
	// Confirmation email customization
	ConfirmUser MailTemplate `json:"confirm_user" yaml:"confirm_user" mapstructure:"confirm_user"`
	// InviteUser email customization
	InviteUser MailTemplate `json:"invite_user" yaml:"invite_user" mapstructure:"invite_user"`
	// ResetPassword email customization
	ResetPassword MailTemplate `json:"reset_password" yaml:"reset_password" mapstructure:"reset_password"`
	// SignupCode email customization
	SignupCode MailTemplate `json:"signupcode"`
}

func (mt *MailTemplates) normalize(srv Service) error {
	referralURLs := []*string{
		&mt.ChangeEmail.ReferralURL,
		&mt.ConfirmUser.ReferralURL,
		&mt.InviteUser.ReferralURL,
		&mt.ResetPassword.ReferralURL,
		&mt.SignupCode.ReferralURL,
	}
	for _, ref := range referralURLs {
		if *ref == "" {
			*ref = srv.SiteURL
		}
		_, err := url.Parse(*ref)
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	mailLinkAction = ":action"
	mailLinkToken  = ":token"
)

const defaultLinkFormat = "/" + mailLinkAction + "/#/" + mailLinkToken

var templateDefaults = MailTemplates{
	Theme: mailTheme,
	ChangeEmail: MailTemplate{
		LinkFormat: defaultLinkFormat,
	},
	ConfirmUser: MailTemplate{
		LinkFormat: defaultLinkFormat,
	},
	InviteUser: MailTemplate{
		LinkFormat: defaultLinkFormat,
	},
	ResetPassword: MailTemplate{
		LinkFormat: defaultLinkFormat,
	},
	SignupCode: MailTemplate{
		LinkFormat: "/" + mailLinkAction,
	},
}

// MailTemplate holds the configuration for emails.
type MailTemplate struct {
	Subject string `json:"subject"`
	// LinkFormat is the url format for the email link. If the format includes ':token' and
	// ':action' strings, they will automatically be replaced with their respective values.
	LinkFormat  string `json:"link_format" yaml:"link_format" mapstructure:"link_format"`
	ReferralURL string `json:"referral_url" yaml:"referral_url" mapstructure:"referral_url"`
	Template    string `json:"template"`
}

// FormatLink formats the link URL replacing the ':token' and ':action' in the link format.
// This function assumes you know what you doing and will produce a valid link.
func FormatLink(linkFormat, action, token string) string {
	l := strings.ReplaceAll(linkFormat, mailLinkAction, action)
	l = strings.ReplaceAll(l, mailLinkToken, token)
	return l
}
