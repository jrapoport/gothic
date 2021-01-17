package api

import (
	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/jrapoport/gothic/conf"
)

func (a *API) checkRecaptcha(ipaddress, token string) error {
	if a.config.Recaptcha.Key == "" {
		return nil
	}
	if token == "" {
		return badRequestError("invalid token")
	} else if ipaddress == "" {
		return badRequestError("invalid ip address")
	}
	err := a.validateRecaptcha(ipaddress, token)
	return err
	/*
		if err == nil {
			return nil
		}
		// FIXME: this is an absurd hack for docker on osx
		// 	https://github.com/docker/for-mac/issues/180
		return a.validateRecaptcha(r.RemoteAddr, token)
	*/
}

const recaptchaDebugKey = "RECAPTCHA-DEBUG-KEY"
const recaptchaDebugToken = "SIGNUP-DEBUG-RECAPTCHA"

func (a *API) validateRecaptcha(remoteIp string, token string) error {
	if conf.Debug && a.config.Recaptcha.Key == recaptchaDebugKey {
		if token == recaptchaDebugToken {
			return nil
		} else {
			return badRequestError("bad debug recaptcha: %s", token)
		}
	}
	recaptcha.Init(a.config.Recaptcha.Key)
	if remoteIp == "" {
		return badRequestError("recaptcha error: %s", "invalid remote ip")
	}
	rc, err := recaptcha.Confirm(remoteIp, token)
	if err != nil {
		return badRequestError("recaptcha error: %v", err)
	}
	if !rc {
		return badRequestError("recaptcha failed")
	}
	return nil
}
