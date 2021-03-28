package provider

import "sync"

const (
	// Unknown provider
	Unknown Name = ""
	// Amazon provider
	Amazon Name = "amazon"
	// Apple provider
	Apple Name = "apple"
	// Auth0 provider
	Auth0 Name = "auth0"
	// AzureAD provider
	AzureAD Name = "azuread"
	// AzureADv2 provider
	AzureADv2 Name = "azureadv2"
	// BattleNet provider
	BattleNet Name = "battlenet"
	// BitBucket provider
	BitBucket Name = "bitbucket"
	// Box provider
	Box Name = "box"
	// CloudFoundry provider
	CloudFoundry Name = "cloudfoundry"
	// DailyMotion provider
	DailyMotion Name = "dailymotion"
	// Deezer provider
	Deezer Name = "deezer"
	// DigitalOcean provider
	DigitalOcean Name = "digitalocean"
	// Discord provider
	Discord Name = "discord"
	// Dropbox provider
	Dropbox Name = "dropbox"
	// EveOnline provider
	EveOnline Name = "eveonline"
	// Facebook provider
	Facebook Name = "facebook"
	// Fitbit provider
	Fitbit Name = "fitbit"
	// GitHub provider
	GitHub Name = "github"
	// GitLab provider
	GitLab Name = "gitlab"
	// Gitea provider
	Gitea Name = "gitea"
	// Google provider
	Google Name = "google"
	// Heroku provider
	Heroku Name = "heroku"
	// InfluxCloud provider
	InfluxCloud Name = "influxcloud"
	// Instagram provider
	Instagram Name = "instagram"
	// Intercom provider
	Intercom Name = "intercom"
	// KaKao provider
	KaKao Name = "kakao"
	// LastFM provider
	LastFM Name = "lastfm"
	// Line provider
	Line Name = "line"
	// LinkedIN provider
	LinkedIN Name = "linkedin"
	// MailRU provider
	MailRU Name = "mailru"
	// Mastodon provider
	Mastodon Name = "mastodon"
	// Meetup provider
	Meetup Name = "meetup"
	// MicrosoftOnline provider
	MicrosoftOnline Name = "microsoftonline"
	// Naver provider
	Naver Name = "naver"
	// NextCloud provider
	NextCloud Name = "nextcloud"
	// Okta provider
	Okta Name = "okta"
	// OneDrive provider
	OneDrive Name = "onedrive"
	// OpenIDConnect provider
	OpenIDConnect Name = "openid-connect"
	// Oura provider
	Oura Name = "oura"
	// PayPal provider
	PayPal Name = "paypal"
	// SalesForce provider
	SalesForce Name = "salesforce"
	// SeaTalk provider
	SeaTalk Name = "seatalk"
	// Shopify provider
	Shopify Name = "shopify"
	// Slack provider
	Slack Name = "slack"
	// SoundCloud provider
	SoundCloud Name = "soundcloud"
	// Spotify provider
	Spotify Name = "spotify"
	// Steam provider
	Steam Name = "steam"
	// Strava provider
	Strava Name = "strava"
	// Stripe provider
	Stripe Name = "stripe"
	// Tumblr provider
	Tumblr Name = "tumblr"
	// Twitch provider
	Twitch Name = "twitch"
	// Twitter provider
	Twitter Name = "twitter"
	// TypeTalk provider
	TypeTalk Name = "typetalk"
	// Uber provider
	Uber Name = "uber"
	// VK provider
	VK Name = "vk"
	// WePay provider
	WePay Name = "wepay"
	// Xero provider
	Xero Name = "xero"
	// Yahoo provider
	Yahoo Name = "yahoo"
	// Yammer provider
	Yammer Name = "yammer"
	// Yandex provider
	Yandex Name = "yandex"
)

// External is a map of supported providers
var External = map[Name]struct{}{
	Amazon: {},
	// FIXME: core/auth/providers.go:162
	// Apple:           {},
	Auth0:           {},
	AzureAD:         {},
	AzureADv2:       {},
	BattleNet:       {},
	BitBucket:       {},
	Box:             {},
	CloudFoundry:    {},
	DailyMotion:     {},
	Deezer:          {},
	DigitalOcean:    {},
	Discord:         {},
	Dropbox:         {},
	EveOnline:       {},
	Facebook:        {},
	Fitbit:          {},
	GitHub:          {},
	GitLab:          {},
	Gitea:           {},
	Google:          {},
	Heroku:          {},
	InfluxCloud:     {},
	Instagram:       {},
	Intercom:        {},
	KaKao:           {},
	LastFM:          {},
	Line:            {},
	LinkedIN:        {},
	MailRU:          {},
	Mastodon:        {},
	Meetup:          {},
	MicrosoftOnline: {},
	Naver:           {},
	NextCloud:       {},
	Okta:            {},
	OneDrive:        {},
	OpenIDConnect:   {},
	Oura:            {},
	PayPal:          {},
	SalesForce:      {},
	SeaTalk:         {},
	Shopify:         {},
	Slack:           {},
	SoundCloud:      {},
	Spotify:         {},
	Steam:           {},
	Strava:          {},
	Stripe:          {},
	Tumblr:          {},
	Twitch:          {},
	Twitter:         {},
	TypeTalk:        {},
	Uber:            {},
	VK:              {},
	WePay:           {},
	Xero:            {},
	Yahoo:           {},
	Yammer:          {},
	Yandex:          {},
}

var mu sync.RWMutex

// AddExternal adds an external provider
func AddExternal(p Name) {
	mu.Lock()
	defer mu.Unlock()
	External[p] = struct{}{}
}

// IsExternal returns true if a provider is external
func IsExternal(p Name) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := External[p]
	return ok
}
