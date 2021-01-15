package storage

import (
	"context"

	"github.com/jrapoport/gothic/conf"
)

const (
	configCtxKey = "global_config"
)

func (c *Connection) withContext(ctx context.Context, global *conf.Configuration) *Connection {
	ctx = withConfig(ctx, global)
	return &Connection{DB: c.DB.WithContext(ctx)}
}

// withConfig adds the tenant configuration to the context.
func withConfig(ctx context.Context, config *conf.Configuration) context.Context {
	return context.WithValue(ctx, configCtxKey, config)
}

// Config gets configuration that was used to initialize the database context
func (c *Connection) Config() *conf.Configuration {
	obj, found := c.Get(configCtxKey)
	if !found {
		return nil
	}
	return obj.(*conf.Configuration)
}
