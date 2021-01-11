package test

import (
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/storage"
)

func SetupDBConnection(globalConfig *conf.GlobalConfiguration) (*storage.Connection, error) {
	return storage.Dial(globalConfig)
}
