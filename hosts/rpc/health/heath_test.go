package health

import (
	"context"
	"testing"

	"github.com/jrapoport/gothic/test/tsrv"
	"github.com/stretchr/testify/assert"
)

func TestHealthServer_HealthCheck(t *testing.T) {
	s, _ := tsrv.RPCServer(t, false)
	srv := newHealthServer(s)
	ctx := context.Background()
	res, err := srv.HealthCheck(ctx, nil)
	assert.NoError(t, err)
	test := s.HealthCheck()
	assert.Equal(t, test.Name, res.Name)
	assert.Equal(t, test.Status, res.Status)
	assert.Equal(t, test.Version, res.Version)
}
