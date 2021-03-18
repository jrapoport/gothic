package audit

import (
	"os"

	"github.com/jrapoport/gothic/models/auditlog"
	"github.com/jrapoport/gothic/models/types"
	"github.com/jrapoport/gothic/models/types/key"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
)

// LogStartup lof the service start
func LogStartup(conn *store.Connection, service string) error {
	name, err := os.Hostname()
	if err != nil {
		return err
	}
	ip, err := utils.OutboundIP()
	if err != nil {
		return err
	}
	_, err = CreateLogEntry(nil, conn, auditlog.Startup, user.SystemID, types.Map{
		key.Service:   service,
		key.Hostname:  name,
		key.IPAddress: ip.String(),
	})
	return err
}

// LogShutdown lof the service shutdown
func LogShutdown(conn *store.Connection, service string) error {
	name, err := os.Hostname()
	if err != nil {
		return err
	}
	ip, err := utils.OutboundIP()
	if err != nil {
		return err
	}
	_, err = CreateLogEntry(nil, conn, auditlog.Shutdown, user.SystemID, types.Map{
		key.Service:   service,
		key.Hostname:  name,
		key.IPAddress: ip.String(),
	})
	return err
}
