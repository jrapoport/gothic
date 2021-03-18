package providers

import (
	"errors"
	"fmt"
	"os"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/store/types/provider"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/azureadv2"
	"github.com/markbates/goth/providers/battlenet"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/cloudfoundry"
	"github.com/markbates/goth/providers/dailymotion"
	"github.com/markbates/goth/providers/deezer"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/eveonline"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/fitbit"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/influxcloud"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/kakao"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/line"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/mailru"
	"github.com/markbates/goth/providers/mastodon"
	"github.com/markbates/goth/providers/meetup"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/naver"
	"github.com/markbates/goth/providers/nextcloud"
	"github.com/markbates/goth/providers/okta"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/oura"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/seatalk"
	"github.com/markbates/goth/providers/shopify"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/strava"
	"github.com/markbates/goth/providers/stripe"
	"github.com/markbates/goth/providers/tumblr"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/typetalk"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/vk"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/xero"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
	"github.com/markbates/goth/providers/yandex"
)

var internalProvider = provider.Unknown

type Provider goth.Provider

// Providers is list of known/available providers.
type Providers map[string]Provider

func NewProviders() Providers {
	return Providers{}
}

// LoadProviders loads the configured external providers.
func (providers *Providers) LoadProviders(c *config.Config) error {
	providers.clearProviders()
	internalProvider = c.Provider()
	for name, v := range c.Providers {
		err := providers.useProvider(name, v.ClientKey, v.Secret, v.CallbackURL, v.Scopes...)
		if err != nil {
			err = fmt.Errorf("load %s provider failed: %w", name, err)
			return err
		}
	}
	return nil
}

// UseProviders adds a list of available providers for use with Goth.
// Can be called multiple times. If you pass the same provider more
// than once, the last will be used.
func (providers Providers) UseProviders(viders ...Provider) {
	for _, p := range viders {
		providers[p.Name()] = p
	}
}

// GetProvider returns a previously created provider. If we have not
// been told to use the named provider it will return an error.
func (providers Providers) GetProvider(p provider.Name) (Provider, error) {
	if !p.IsExternal() {
		err := fmt.Errorf("invalid provider: %s", p)
		return nil, err
	}
	return providers.getProvider(p.String())
}

func (providers Providers) getProvider(name string) (Provider, error) {
	p := providers[name]
	if p == nil {
		return nil, fmt.Errorf("no provider for %s exists", name)
	}
	return p, nil
}

// IsEnabled returns true if the provider is enabled.
func (providers Providers) IsEnabled(p provider.Name) error {
	// check against Unknown first so we catch internalProvider
	// properly (in case internalProvider is Unknown)
	if p == provider.Unknown {
		return errors.New("invalid provider")
	} else if p == internalProvider {
		return nil
	}
	_, err := providers.GetProvider(p)
	return err
}

// ClearProviders will remove all providers currently in use.
// This is useful, mostly, for testing purposes.
func (providers *Providers) clearProviders() {
	internalProvider = provider.Unknown
	*providers = Providers{}
}

func (providers *Providers) useProvider(name provider.Name, clientKey, secret, callback string, scopes ...string) (err error) {
	var p Provider
	switch name {
	case provider.Amazon:
		p = amazon.New(clientKey, secret, callback, scopes...)
	case provider.Apple:
		scopes = append(scopes, apple.ScopeName, apple.ScopeEmail)
		p = apple.New(clientKey, secret, callback, nil, scopes...)
	case provider.Auth0:
		// your Auth0 customer domain is required
		domain := getEnv(config.Auth0DomainEnv)
		if domain == "" {
			return errors.New("auth0 customer domain required")
		}
		p = auth0.New(clientKey, secret, callback, domain, scopes...)
	case provider.AzureAD:
		p = azuread.New(clientKey, secret, callback, nil, scopes...)
	case provider.AzureADv2:
		opts := azureadv2.ProviderOptions{}
		opts.Scopes = make([]azureadv2.ScopeType, len(scopes))
		for i, scope := range scopes {
			opts.Scopes[i] = azureadv2.ScopeType(scope)
		}
		if tnt := getEnv(config.AzureADTenantEnv); tnt != "" {
			opts.Tenant = azureadv2.TenantType(tnt)
		}
		p = azureadv2.New(clientKey, secret, callback, opts)
	case provider.BattleNet:
		p = battlenet.New(clientKey, secret, callback, scopes...)
	case provider.BitBucket:
		p = bitbucket.New(clientKey, secret, callback, scopes...)
	case provider.Box:
		p = box.New(clientKey, secret, callback, scopes...)
	case provider.CloudFoundry:
		url := getEnv(config.CloudFoundryURLEnv)
		p = cloudfoundry.New(url, clientKey, secret, callback, scopes...)
	case provider.DailyMotion:
		scopes = append(scopes, "email")
		p = dailymotion.New(clientKey, secret, callback, scopes...)
	case provider.Deezer:
		scopes = append(scopes, "email")
		p = deezer.New(clientKey, secret, callback, scopes...)
	case provider.DigitalOcean:
		scopes = append(scopes, "read")
		p = digitalocean.New(clientKey, secret, callback, scopes...)
	case provider.Discord:
		scopes = append(scopes, discord.ScopeIdentify, discord.ScopeEmail)
		p = discord.New(clientKey, secret, callback, scopes...)
	case provider.Dropbox:
		p = dropbox.New(clientKey, secret, callback, scopes...)
	case provider.EveOnline:
		p = eveonline.New(clientKey, secret, callback, scopes...)
	case provider.Facebook:
		p = facebook.New(clientKey, secret, callback, scopes...)
	case provider.Fitbit:
		p = fitbit.New(clientKey, secret, callback, scopes...)
	case provider.Gitea:
		p = gitea.New(clientKey, secret, callback, scopes...)
	case provider.GitHub:
		scopes = append(scopes, "user:email")
		p = github.New(clientKey, secret, callback, scopes...)
	case provider.GitLab:
		p = gitlab.New(clientKey, secret, callback, scopes...)
	case provider.Google:
		p = google.New(clientKey, secret, callback, scopes...)
	case provider.Heroku:
		p = heroku.New(clientKey, secret, callback, scopes...)
	case provider.InfluxCloud:
		p = influxcloud.New(clientKey, secret, callback, scopes...)
	case provider.Instagram:
		p = instagram.New(clientKey, secret, callback, scopes...)
	case provider.Intercom:
		p = intercom.New(clientKey, secret, callback, scopes...)
	case provider.KaKao:
		p = kakao.New(clientKey, secret, callback, scopes...)
	case provider.LastFM:
		p = lastfm.New(clientKey, secret, callback)
	case provider.Line:
		scopes = append(scopes, "profile", "openid", "email")
		p = line.New(clientKey, secret, callback, scopes...)
	case provider.LinkedIN:
		p = linkedin.New(clientKey, secret, callback, scopes...)
	case provider.MailRU:
		p = mailru.New(clientKey, secret, callback, scopes...)
	case provider.Mastodon:
		scopes = append(scopes, "read:accounts")
		p = mastodon.New(clientKey, secret, callback, scopes...)
	case provider.Meetup:
		p = meetup.New(clientKey, secret, callback, scopes...)
	case provider.MicrosoftOnline:
		p = microsoftonline.New(clientKey, secret, callback, scopes...)
	case provider.Naver:
		p = naver.New(clientKey, secret, callback)
	case provider.NextCloud:
		url := getEnv(config.NextCloudURLEnv)
		p = nextcloud.NewCustomisedDNS(clientKey, secret, callback, url, scopes...)
	case provider.Okta:
		url := getEnv(config.OktaURLEnv)
		scopes = append(scopes, "openid", "profile", "email")
		p = okta.New(clientKey, secret, url, callback, scopes...)
	case provider.OneDrive:
		p = onedrive.New(clientKey, secret, callback, scopes...)
	case provider.OpenIDConnect:
		// auto discovery url (https://openid.net/specs/openid-connect-discovery-1_0-17.html).
		url := getEnv(config.OpenIDConnectURLEnv)
		if url == "" {
			return errors.New("openid connect discovery url required")
		}
		p, err = openidConnect.New(clientKey, secret, callback, url, scopes...)
		if err != nil {
			return err
		}
	case provider.Oura:
		p = oura.New(clientKey, secret, callback, scopes...)
	case provider.PayPal:
		// set PAYPAL_ENV=sandbox as environment variable to use the paypal sandbox.
		p = paypal.New(clientKey, secret, callback, scopes...)
	case provider.SalesForce:
		p = salesforce.New(clientKey, secret, callback, scopes...)
	case provider.SeaTalk:
		p = seatalk.New(clientKey, secret, callback, scopes...)
	case provider.Shopify:
		scopes = append(scopes, shopify.ScopeReadCustomers, shopify.ScopeReadOrders)
		p = shopify.New(clientKey, secret, callback, scopes...)
	case provider.Slack:
		p = slack.New(clientKey, secret, callback, scopes...)
	case provider.SoundCloud:
		p = soundcloud.New(clientKey, secret, callback, scopes...)
	case provider.Spotify:
		p = spotify.New(clientKey, secret, callback, scopes...)
	case provider.Steam:
		p = steam.New(clientKey, callback)
	case provider.Strava:
		p = strava.New(clientKey, secret, callback, scopes...)
	case provider.Stripe:
		p = stripe.New(clientKey, secret, callback, scopes...)
	case provider.Tumblr:
		p = tumblr.New(clientKey, secret, callback)
	case provider.Twitch:
		p = twitch.New(clientKey, secret, callback, scopes...)
	case provider.Twitter:
		if getEnv(config.TwitterAuthorizeEnv) != "" {
			// use authorize instead of authenticate with twitter
			p = twitter.New(clientKey, secret, callback)
		} else {
			p = twitter.NewAuthenticate(clientKey, secret, callback)
		}
	case provider.TypeTalk:
		scopes = append(scopes, "my")
		p = typetalk.New(clientKey, secret, callback, scopes...)
	case provider.Uber:
		p = uber.New(clientKey, secret, callback, scopes...)
	case provider.VK:
		p = vk.New(clientKey, secret, callback, scopes...)
	case provider.WePay:
		scopes = append(scopes, "view_user")
		p = wepay.New(clientKey, secret, callback, scopes...)
	case provider.Xero:
		p = xero.New(clientKey, secret, callback)
	case provider.Yahoo:
		// pointed localhost.com to http://localhost:3000/auth/yahoo/callback through proxy as yahoo
		// does not allow to put custom ports in redirection uri
		p = yahoo.New(clientKey, secret, "http://localhost.com", scopes...)
	case provider.Yammer:
		p = yammer.New(clientKey, secret, callback, scopes...)
	case provider.Yandex:
		p = yandex.New(clientKey, secret, callback, scopes...)
	default:
		return fmt.Errorf("invalid provider: %s", name)
	}
	providers.UseProviders(p)
	return nil
}

func getEnv(key string) string {
	return os.Getenv(config.ENVPrefix + "_" + key)
}
