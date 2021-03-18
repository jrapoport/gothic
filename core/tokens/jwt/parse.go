package jwt

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types"
)

// ParseClaims parses a set of jwt claims from a token.
func ParseClaims(c config.JWT, token string, claims Claims) error {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		clm, ok := claims.(Claims)
		if !ok {
			return nil, errors.New("invalid claims")
		}
		sc := clm.Standard()
		if sc.Subject == "" {
			return nil, errors.New("invalid subject")
		}
		return []byte(c.Secret), nil
	}
	var opts []jwt.ParserOption
	if c.Issuer != "" {
		opt := jwt.WithIssuer(c.Issuer)
		opts = append(opts, opt)
	}
	if c.Audience != "" {
		auds := strings.Split(c.Audience, ",")
		for _, aud := range auds {
			opt := jwt.WithAudience(aud)
			opts = append(opts, opt)
		}
	} else {
		opt := jwt.WithoutAudienceValidation()
		opts = append(opts, opt)
	}
	_, err := jwt.ParseWithClaims(token, claims, keyFunc, opts...)
	return err
}

// ParseData parses a Map from a token.
func ParseData(c config.JWT, token string) (types.Map, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(c.Secret), nil
	}
	t, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return nil, err
	}
	return t.Header, nil
}

// TODO: support RSA & ECDSA keys.
// 	jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) {
//		 return LoadRSAPublicKeyFromDisk("test/sample_key.pub"), nil
//		 }
