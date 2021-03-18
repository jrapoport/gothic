package audit

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
	"github.com/stretchr/testify/require"
)

func TestLogStartup(t *testing.T) {
	t.Parallel()
	const service = "test_service"
	host, err := os.Hostname()
	require.NoError(t, err)
	ipaddr, err := utils.OutboundIP()
	require.NoError(t, err)
	testLogEntry(t, auditlog.Startup, uuid.Nil,
		types.Map{
			key.Service:   service,
			key.IPAddress: ipaddr.String(),
			key.Hostname:  host,
		},
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogStartup(conn, service)
		})
}

func TestLogShutdown(t *testing.T) {
	t.Parallel()
	const service = "test_service"
	host, err := os.Hostname()
	require.NoError(t, err)
	ipaddr, err := utils.OutboundIP()
	require.NoError(t, err)
	testLogEntry(t, auditlog.Shutdown, uuid.Nil,
		types.Map{
			key.Service:   service,
			key.IPAddress: ipaddr.String(),
			key.Hostname:  host,
		},
		func(ctx context.Context, conn *store.Connection, uid uuid.UUID, _ types.Map) error {
			return LogShutdown(conn, service)
		})
}
