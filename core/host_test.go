package core

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/jrapoport/gothic/test/tconf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHost(t *testing.T) {
	t.Parallel()
	const (
		address = "127.0.0.1:0"
		name    = "test"
	)
	c := tconf.TempDB(t)
	a, err := NewAPI(c)
	require.NoError(t, err)
	h := NewHost(a, name, "")
	assert.NotNil(t, h)
	assert.Empty(t, h.Address())
	err = h.ListenAndServe()
	assert.NoError(t, err)
	assert.False(t, h.Online())
	h = NewHost(a, name, address)
	h.Start(func(l net.Listener) error {
		_, err = l.Accept()
		assert.Error(t, err)
		return err
	})
	err = h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.lis != nil
	}, 1*time.Second, 10*time.Millisecond)
	assert.True(t, h.Online())
	assert.Equal(t, name, h.Name())
	assert.NotEmpty(t, h.Address())
	assert.NotEqual(t, address, h.Address())
	h.Stop(func(ctx context.Context) error {
		return nil
	})
	err = h.Shutdown()
	assert.NoError(t, err)
}

func TestNewHost_Error(t *testing.T) {
	t.Parallel()
	c := tconf.TempDB(t)
	a, err := NewAPI(c)
	require.NoError(t, err)
	h := NewHost(a, "test", "255.255.255.255")
	h.Start(func(l net.Listener) error {
		return nil
	})
	err = h.ListenAndServe()
	assert.Error(t, err)
	h = NewHost(a, "test", "127.0.0.1:0")
	h.Start(func(l net.Listener) error {
		_, err = l.Accept()
		assert.Error(t, err)
		return err
	})
	err = h.ListenAndServe()
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		return h.lis != nil
	}, 1*time.Second, 10*time.Millisecond)
	h.Stop(func(ctx context.Context) error {
		return errors.New("fake")
	})
	err = h.Shutdown()
	assert.Error(t, err)
}
