package config

import (
	"errors"
	"fmt"
	"github.com/jrapoport/gothic/utils"
	"strings"
	"time"

	"github.com/imdario/mergo"
)

// JWT holds all the JWT related configuration.
type JWT struct {
	Secret     string `json:"secret"`
	PEM        `yaml:",inline" mapstructure:",squash"`
	Algorithm  string        `json:"algorithm"`
	Expiration time.Duration `json:"expiration"`
	// Issuer is the the entity that issued the token (default: Config.Service)
	Issuer string `json:"issuer"`
	// Audience is an optional comma separated list of resource
	// servers that should accept the token (default: n/a)
	Audience string `json:"audience"`
	// Scope is an optional comma separated list of permission
	// scopes for the token (default: n/a)
	Scope string `json:"scope"`
}

type PEM struct {
	PrivateKey string `json:"privatekey"`
	PublicKey  string `json:"publickey"`
}

func (j *JWT) normalize(srv Service, def JWT) {
	if def.Issuer == "" {
		def.Issuer = strings.ToLower(srv.Name)
	}
	// no error is possible here since we
	// control the struct entirely
	_ = mergo.Merge(j, def)
}

func (j *JWT) CheckRequired() error {
	if j.Secret == "" && j.PrivateKey == "" {
		return errors.New("jwt secret or private key is required")
	} else if j.Secret != "" {
		return nil
	}
	if j.PrivateKey != "" && !utils.PathExists(j.PrivateKey) {
		err := fmt.Errorf("jwt private key not found: %s", j.PrivateKey)
		return err
	}
	if j.PublicKey != "" && !utils.PathExists(j.PublicKey) {
		err := fmt.Errorf("jwt public key not found: %s", j.PublicKey)
		return err
	}
	return nil
}
