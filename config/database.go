package config

import (
	"errors"

	"github.com/jrapoport/gothic/store/drivers"
	"github.com/jrapoport/gothic/utils"
)

// Database holds all the database related configuration.
type Database struct {
	Namespace   string         `json:"namespace"`
	Driver      drivers.Driver `json:"driver"`
	DSN         string         `json:"dsn"`
	MaxRetries  int            `json:"max_retries" yaml:"max_retries" mapstructure:"max_retries"`
	AutoMigrate bool           `json:"automigrate"`
}

func (d *Database) normalize(srv Service) (err error) {
	if d.Driver == "" {
		return errors.New("database driver required")
	}
	if d.DSN == "" {
		return errors.New("database dsn required")
	}
	name := utils.Namespaced(d.Namespace, srv.Name)
	_, d.DSN, err = drivers.NormalizeDSN(name, d.Driver, d.DSN)
	return
}
