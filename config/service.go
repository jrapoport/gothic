package config

import (
	"errors"
	"net/url"
	"strings"
)

// Service config
type Service struct {
	// Service the name of the microservice. default: gothic
	Name string `json:"service" yaml:"service" mapstructure:"service"`
	// SiteURL is the url for the microservice.
	SiteURL string `json:"site_url" yaml:"site_url" mapstructure:"site_url"`
}

// Version returns the build version.
func (s Service) Version() string {
	return BuildVersion()
}

func (s *Service) normalize() error {
	// make sure this wasn't mucked with
	uri, err := url.Parse(s.SiteURL)
	if err != nil {
		return err
	}
	if s.Name == "" {
		s.Name = strings.ToLower(uri.Host)
	}
	return nil
}

func (s *Service) CheckRequired() error {
	if s.SiteURL == "" {
		return errors.New("site url is required")
	}
	return nil
}
