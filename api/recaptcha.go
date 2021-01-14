package api

import (
	"net"
	"net/http"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/jrapoport/gothic/conf"
	"github.com/sirupsen/logrus"
)

func (a *API) checkRecaptcha(r *http.Request, config *conf.Configuration) error {
	if config.Recaptcha.Key == "" {
		return nil
	}
	recap := r.FormValue("recaptcha")
	if recap == "" {
		return badRequestError("invalid recaptcha")
	}
	if err := a.validateRecaptcha(a.getUserAddress(r), recap); err != nil {
		// FIXME: this is an absurd hack for docker on osx
		// 	https://github.com/docker/for-mac/issues/180
		if err = a.validateRecaptcha(r.RemoteAddr, recap); err != nil {
			return err
		}
	}
	return nil
}

func (a *API) validateRecaptcha(remoteIp string, token string) error {
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

func (a *API) getUserAddress(r *http.Request) string {
	addr := r.Header.Get("X-Real-Ip")
	logrus.Infof("X-Real-Ip: %s", addr)
	if addr == "" {
		addr = r.Header.Get("X-Forwarded-For")
		logrus.Infof("X-Forwarded-For: %s", addr)
	}
	if addr == "" {
		addr = r.RemoteAddr
		logrus.Infof("RemoteAddr: %s", addr)
	}
	addr, _, err := net.SplitHostPort(addr)
	if err != nil {
		return ""
	}
	return addr
}
