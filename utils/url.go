package utils

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/purell"
	"github.com/segmentio/encoding/json"
)

// NormalizeURL returns a normalized a url.
func NormalizeURL(u string) (string, error) {
	const flags = purell.FlagsSafe | purell.FlagRemoveDotSegments |
		purell.FlagRemoveDirectoryIndex | purell.FlagRemoveDuplicateSlashes |
		purell.FlagRemoveUnnecessaryHostDots | purell.FlagRemoveEmptyPortSeparator
	return purell.NormalizeURLString(u, flags)
}

// JoinLink returns a normalized url with the attached fragment and base.
func JoinLink(linkURL, fragment string) (string, error) {
	base, err := url.Parse(linkURL)
	if err != nil {
		return "", err
	}
	link, err := appendFragment(base, fragment)
	if err != nil {
		return "", err
	}
	return NormalizeURL(link.String())
}

func appendFragment(base *url.URL, fragmentURL string) (*url.URL, error) {
	// make sure this is fragment
	frag, err := url.Parse(makeRelative(fragmentURL))
	if err != nil {
		return nil, err
	}
	return base.ResolveReference(frag), nil
}

var relRx = regexp.MustCompile(`^https?://[^/]+`)

func makeRelative(url string) string {
	url = relRx.ReplaceAllString(url, "")
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	return url
}

// URLValuesToMap converts a map to url.Values with json data support
func URLValuesToMap(values url.Values, mapDataKey bool) map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range values {
		if len(v) <= 0 {
			continue
		}
		m[k] = v[0]
		if mapDataKey && k == "data" {
			data := map[string]interface{}{}
			err := json.Unmarshal([]byte(v[0]), &data)
			if err == nil {
				m[k] = data
			}
		}
	}
	return m
}
