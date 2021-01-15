package conf

import "github.com/dgrijalva/jwt-go/v4"

// JWTConfig holds all the JWT related configuration.
type JWTConfig struct {
	Secret       string `json:"-" required:"true"`
	Method       string `json:"method" default:"HS256"`
	Subject      string `json:"subject" default:"gothic"`
	Exp          int    `json:"exp"`
	Aud          string `json:"aud"`
	AdminGroup   string `json:"admin_group" split_words:"true" default:"admin"`
	DefaultGroup string `json:"default_group" split_words:"true" default:"user"`
	MaskEmail    bool   `json:"mask_email" split_words:"true"`
}

func (c *JWTConfig) SigningMethod() jwt.SigningMethod {
	return jwt.GetSigningMethod(c.Method)
}
