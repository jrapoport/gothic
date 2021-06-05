package jwt

import (
	"errors"
	"strings"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/models/types"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// ParseClaims parses a set of jwt claims from a token.
func ParseClaims(c config.JWT, token string, claims Claims) error {
	tok, err := parse(c, token)
	if err != nil {
		return err
	}
	sub := tok.Subject()
	if sub == "" {
		return errors.New("invalid subject")
	}
	claims.parseToken(&Token{Token: tok})
	return nil
}

// ParseData parses a Map from a token.
func ParseData(c config.JWT, token string) (types.Map, error) {
	tok, err := parse(c, token)
	if err != nil {
		return nil, err
	}
	return tok.PrivateClaims(), nil
}

func parse(c config.JWT, token string) (jwt.Token, error) {
	alg := jwa.SignatureAlgorithm(c.Algorithm)
	key := publicKey(c)
	opts := []jwt.ParseOption{
		jwt.WithValidate(true),
		jwt.WithVerify(alg, key),
	}
	if c.Issuer != "" {
		opt := jwt.WithIssuer(c.Issuer)
		opts = append(opts, opt)
	}
	if c.Audience != "" {
		aud := strings.Split(c.Audience, ",")
		for _, a := range aud {
			opt := jwt.WithAudience(a)
			opts = append(opts, opt)
		}
	}
	tok, err := jwt.Parse([]byte(token), opts...)
	if err != nil {
		return nil, err
	}
	return tok, nil
}
