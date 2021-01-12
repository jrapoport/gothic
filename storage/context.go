package storage

import (
	"context"

	"github.com/jrapoport/gothic/conf"
)

const (
	globalConfigCtxKey = "global_config"
)

func (c *Connection) withContext(ctx context.Context, global *conf.GlobalConfiguration) *Connection {
	ctx = withGlobalConfig(ctx, global)
	return &Connection{DB: c.DB.WithContext(ctx)}
}

// withConfig adds the tenant configuration to the context.
func withGlobalConfig(ctx context.Context, config *conf.GlobalConfiguration) context.Context {
	return context.WithValue(ctx, globalConfigCtxKey, config)
}

func (c *Connection) GetGlobalConfig(ctx context.Context) *conf.GlobalConfiguration {
	obj, found := c.Get(globalConfigCtxKey)
	if !found {
		return nil
	}
	return obj.(*conf.GlobalConfiguration)
}
