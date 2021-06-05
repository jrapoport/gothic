package jwt

import (
	"io/ioutil"

	"github.com/jrapoport/gothic/config"
	"github.com/lestrrat-go/jwx/jwk"
)

func privateKey(j config.JWT) (key jwk.Key) {
	if j.Secret != "" {
		sec := []byte(j.Secret)
		key, _ = jwk.New(sec)
	} else {
		key, _ = pemKey(j.PEM.PrivateKey)
	}
	return
}

func publicKey(j config.JWT) (key jwk.Key) {
	if j.Secret != "" {
		sec := []byte(j.Secret)
		key, _ = jwk.New(sec)
	} else {
		pem := j.PEM.PublicKey
		if pem == "" {
			const publicKeyExt = ".pub"
			pem = j.PEM.PrivateKey + publicKeyExt
		}
		key, _ = pemKey(pem)
	}
	return
}

func pemKey(pem string) (jwk.Key, error) {
	raw, err := ioutil.ReadFile(pem)
	if err != nil {
		return nil, err
	}
	key, err := jwk.ParseKey(raw, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}
	return key, nil
}
