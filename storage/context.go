package storage

import (
	"context"

	"github.com/jrapoport/gothic/conf"
)

const (
	globalConfigCtxKey = "global_config"
)

func (c *Connection) withContext(ctx context.Context, global *conf.Configuration) *Connection {
	ctx = withGlobalConfig(ctx, global)
	return &Connection{DB: c.DB.WithContext(ctx)}
}

// withConfig adds the tenant configuration to the context.
func withGlobalConfig(ctx context.Context, config *conf.Configuration) context.Context {
	return context.WithValue(ctx, globalConfigCtxKey, config)
}

func (c *Connection) GetGlobalConfig(ctx context.Context) *conf.Configuration {
	obj, found := c.Get(globalConfigCtxKey)
	if !found {
		return nil
	}
	return obj.(*conf.Configuration)
}
