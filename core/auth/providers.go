package auth

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types/provider"
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

// Provider wraps a Provider
type Provider goth.Provider

// Providers is list of known/available providers.
type Providers struct {
	internal  provider.Name
	providers map[string]Provider
	mu        sync.RWMutex
}

// NewProviders returns a new providers
func NewProviders() *Providers {
	return &Providers{
		internal:  provider.Unknown,
		providers: map[string]Provider{},
	}
}

// LoadProviders loads the configured external providers.
func (pv *Providers) LoadProviders(c *config.Config) error {
	pv.mu.Lock()
	defer pv.mu.Unlock()
	pv.internal = c.Provider()
	for name, v := range c.Providers {
		err := pv.useProvider(name, v.ClientKey, v.Secret, v.CallbackURL, v.Scopes...)
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
func (pv *Providers) UseProviders(ps ...Provider) {
	pv.mu.Lock()
	defer pv.mu.Unlock()
	for _, p := range ps {
		pv.providers[p.Name()] = p
	}
}

// GetProvider returns a previously created provider. If we have not
// been told to use the named provider it will return an error.
func (pv *Providers) GetProvider(name provider.Name) (Provider, error) {
	pv.mu.RLock()
	defer pv.mu.RUnlock()
	// properly (in case internalProvider is Unknown)
	if name == provider.Unknown {
		return nil, errors.New("invalid provider")
	}
	if !name.IsExternal() {
		err := fmt.Errorf("invalid provider: %s", name)
		return nil, err
	}
	pvd := pv.providers[string(name)]
	if pvd == nil {
		return nil, fmt.Errorf("no provider for %s exists", name)
	}
	return pvd, nil
}

// IsEnabled returns true if the provider is enabled.
func (pv *Providers) IsEnabled(name provider.Name) error {
	if pv.isInternal(name) {
		return nil
	}
	_, err := pv.GetProvider(name)
	return err
}

func (pv *Providers) isInternal(name provider.Name) bool {
	pv.mu.RLock()
	defer pv.mu.RUnlock()
	return name == pv.internal && pv.internal != provider.Unknown
}

func (pv *Providers) useProvider(name provider.Name, clientKey, secret, callback string, scopes ...string) (err error) {
	var pvdr Provider
	switch name {
	case provider.Amazon:
		pvdr = amazon.New(clientKey, secret, callback, scopes...)
	case provider.Apple:
		scopes = append(scopes, apple.ScopeName, apple.ScopeEmail)
		pvdr = apple.New(clientKey, secret, callback, nil, scopes...)
	case provider.Auth0:
		// your Auth0 customer domain is required
		domain := getEnv(config.Auth0DomainEnv)
		if domain == "" {
			return errors.New("auth0 customer domain required")
		}
		pvdr = auth0.New(clientKey, secret, callback, domain, scopes...)
	case provider.AzureAD:
		pvdr = azuread.New(clientKey, secret, callback, nil, scopes...)
	case provider.AzureADv2:
		opts := azureadv2.ProviderOptions{}
		opts.Scopes = make([]azureadv2.ScopeType, len(scopes))
		for i, scope := range scopes {
			opts.Scopes[i] = azureadv2.ScopeType(scope)
		}
		if tnt := getEnv(config.AzureADTenantEnv); tnt != "" {
			opts.Tenant = azureadv2.TenantType(tnt)
		}
		pvdr = azureadv2.New(clientKey, secret, callback, opts)
	case provider.BattleNet:
		pvdr = battlenet.New(clientKey, secret, callback, scopes...)
	case provider.BitBucket:
		pvdr = bitbucket.New(clientKey, secret, callback, scopes...)
	case provider.Box:
		pvdr = box.New(clientKey, secret, callback, scopes...)
	case provider.CloudFoundry:
		url := getEnv(config.CloudFoundryURLEnv)
		pvdr = cloudfoundry.New(url, clientKey, secret, callback, scopes...)
	case provider.DailyMotion:
		scopes = append(scopes, "email")
		pvdr = dailymotion.New(clientKey, secret, callback, scopes...)
	case provider.Deezer:
		scopes = append(scopes, "email")
		pvdr = deezer.New(clientKey, secret, callback, scopes...)
	case provider.DigitalOcean:
		scopes = append(scopes, "read")
		pvdr = digitalocean.New(clientKey, secret, callback, scopes...)
	case provider.Discord:
		scopes = append(scopes, discord.ScopeIdentify, discord.ScopeEmail)
		pvdr = discord.New(clientKey, secret, callback, scopes...)
	case provider.Dropbox:
		pvdr = dropbox.New(clientKey, secret, callback, scopes...)
	case provider.EveOnline:
		pvdr = eveonline.New(clientKey, secret, callback, scopes...)
	case provider.Facebook:
		pvdr = facebook.New(clientKey, secret, callback, scopes...)
	case provider.Fitbit:
		pvdr = fitbit.New(clientKey, secret, callback, scopes...)
	case provider.Gitea:
		pvdr = gitea.New(clientKey, secret, callback, scopes...)
	case provider.GitHub:
		scopes = append(scopes, "user:email")
		pvdr = github.New(clientKey, secret, callback, scopes...)
	case provider.GitLab:
		pvdr = gitlab.New(clientKey, secret, callback, scopes...)
	case provider.Google:
		pvdr = google.New(clientKey, secret, callback, scopes...)
	case provider.Heroku:
		pvdr = heroku.New(clientKey, secret, callback, scopes...)
	case provider.InfluxCloud:
		pvdr = influxcloud.New(clientKey, secret, callback, scopes...)
	case provider.Instagram:
		pvdr = instagram.New(clientKey, secret, callback, scopes...)
	case provider.Intercom:
		pvdr = intercom.New(clientKey, secret, callback, scopes...)
	case provider.KaKao:
		pvdr = kakao.New(clientKey, secret, callback, scopes...)
	case provider.LastFM:
		pvdr = lastfm.New(clientKey, secret, callback)
	case provider.Line:
		scopes = append(scopes, "profile", "openid", "email")
		pvdr = line.New(clientKey, secret, callback, scopes...)
	case provider.LinkedIN:
		pvdr = linkedin.New(clientKey, secret, callback, scopes...)
	case provider.MailRU:
		pvdr = mailru.New(clientKey, secret, callback, scopes...)
	case provider.Mastodon:
		scopes = append(scopes, "read:accounts")
		pvdr = mastodon.New(clientKey, secret, callback, scopes...)
	case provider.Meetup:
		pvdr = meetup.New(clientKey, secret, callback, scopes...)
	case provider.MicrosoftOnline:
		pvdr = microsoftonline.New(clientKey, secret, callback, scopes...)
	case provider.Naver:
		pvdr = naver.New(clientKey, secret, callback)
	case provider.NextCloud:
		url := getEnv(config.NextCloudURLEnv)
		pvdr = nextcloud.NewCustomisedDNS(clientKey, secret, callback, url, scopes...)
	case provider.Okta:
		url := getEnv(config.OktaURLEnv)
		scopes = append(scopes, "openid", "profile", "email")
		pvdr = okta.New(clientKey, secret, url, callback, scopes...)
	case provider.OneDrive:
		pvdr = onedrive.New(clientKey, secret, callback, scopes...)
	case provider.OpenIDConnect:
		// auto discovery url (https://openid.net/specs/openid-connect-discovery-1_0-17.html).
		url := getEnv(config.OpenIDConnectURLEnv)
		if url == "" {
			return errors.New("openid connect discovery url required")
		}
		pvdr, err = openidConnect.New(clientKey, secret, callback, url, scopes...)
		if err != nil {
			return err
		}
	case provider.Oura:
		pvdr = oura.New(clientKey, secret, callback, scopes...)
	case provider.PayPal:
		// set PAYPAL_ENV=sandbox as environment variable to use the paypal sandbox.
		pvdr = paypal.New(clientKey, secret, callback, scopes...)
	case provider.SalesForce:
		pvdr = salesforce.New(clientKey, secret, callback, scopes...)
	case provider.SeaTalk:
		pvdr = seatalk.New(clientKey, secret, callback, scopes...)
	case provider.Shopify:
		scopes = append(scopes, shopify.ScopeReadCustomers, shopify.ScopeReadOrders)
		pvdr = shopify.New(clientKey, secret, callback, scopes...)
	case provider.Slack:
		pvdr = slack.New(clientKey, secret, callback, scopes...)
	case provider.SoundCloud:
		pvdr = soundcloud.New(clientKey, secret, callback, scopes...)
	case provider.Spotify:
		pvdr = spotify.New(clientKey, secret, callback, scopes...)
	case provider.Steam:
		pvdr = steam.New(clientKey, callback)
	case provider.Strava:
		pvdr = strava.New(clientKey, secret, callback, scopes...)
	case provider.Stripe:
		pvdr = stripe.New(clientKey, secret, callback, scopes...)
	case provider.Tumblr:
		pvdr = tumblr.New(clientKey, secret, callback)
	case provider.Twitch:
		pvdr = twitch.New(clientKey, secret, callback, scopes...)
	case provider.Twitter:
		if getEnv(config.TwitterAuthorizeEnv) != "" {
			// use authorize instead of authenticate with twitter
			pvdr = twitter.New(clientKey, secret, callback)
		} else {
			pvdr = twitter.NewAuthenticate(clientKey, secret, callback)
		}
	case provider.TypeTalk:
		scopes = append(scopes, "my")
		pvdr = typetalk.New(clientKey, secret, callback, scopes...)
	case provider.Uber:
		pvdr = uber.New(clientKey, secret, callback, scopes...)
	case provider.VK:
		pvdr = vk.New(clientKey, secret, callback, scopes...)
	case provider.WePay:
		scopes = append(scopes, "view_user")
		pvdr = wepay.New(clientKey, secret, callback, scopes...)
	case provider.Xero:
		pvdr = xero.New(clientKey, secret, callback)
	case provider.Yahoo:
		// pointed localhost.com to http://localhost:3000/auth/yahoo/callback through proxy as yahoo
		// does not allow to put custom ports in redirection uri
		pvdr = yahoo.New(clientKey, secret, "http://localhost.com", scopes...)
	case provider.Yammer:
		pvdr = yammer.New(clientKey, secret, callback, scopes...)
	case provider.Yandex:
		pvdr = yandex.New(clientKey, secret, callback, scopes...)
	default:
		return fmt.Errorf("invalid provider: %s", name)
	}
	pv.providers[pvdr.Name()] = pvdr
	return nil
}

func getEnv(key string) string {
	return os.Getenv(config.ENVPrefix + "_" + key)
}
