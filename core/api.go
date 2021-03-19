package core

import (
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core/audit"
	"github.com/jrapoport/gothic/core/events"
	"github.com/jrapoport/gothic/core/providers"
	"github.com/jrapoport/gothic/mail"
	"github.com/jrapoport/gothic/models/types/provider"
	"github.com/jrapoport/gothic/store"
	"github.com/sirupsen/logrus"
)

// API is the main API
type API struct {
	config *config.Config
	conn   *store.Connection
	evt    *events.Dispatch
	mail   *mail.Client
	ext    *providers.Providers
	log    logrus.FieldLogger
}

// NewAPI creates a new core API with a configured storage connection
func NewAPI(c *config.Config) (*API, error) {
	a := new(API)
	err := a.LoadConfig(c)
	if err != nil {
		c.Log().Error(err)
		return nil, err
	}
	return a, nil
}

// LoadConfig loads the config
func (a *API) LoadConfig(c *config.Config) (err error) {
	a.config = c
	l := a.config.Log()
	if l == nil {
		l = logrus.New()
	}
	// set the log first so we can log other errors appropriately
	a.log = l.WithField("api", a.config.Env())
	a.evt = events.NewDispatch(c.Name, l)
	a.conn, err = store.Dial(c, a.log)
	if err != nil {
		return a.logError(err)
	}
	err = a.conn.AutoMigrate()
	if err != nil {
		return a.logError(err)
	}
	err = a.CreateSuperAdmin()
	if err != nil {
		return a.logError(err)
	}
	err = a.OpenMail()
	if err != nil {
		return a.logError(err)
	}
	a.ext = providers.NewProviders()
	err = a.ext.LoadProviders(a.config)
	if err != nil {
		return a.logError(err)
	}
	err = audit.LogStartup(a.conn, a.config.Service.Name)
	return a.logError(err)
}

// Providers is mainly here for testing rn
func (a *API) Providers() *providers.Providers {
	return a.ext
}

// Provider returns the name of the internal provider.
func (a *API) Provider() provider.Name {
	return a.config.Provider()
}

// Shutdown shuts down the api service
func (a *API) Shutdown() error {
	a.CloseMail()
	a.closeDispatch()
	err := audit.LogShutdown(a.conn, a.config.Service.Name)
	return a.logError(err)
}

func (a *API) logError(err error) error {
	if err == nil {
		return nil
	}
	a.log.Error(err)
	return err
}
