package cmd

import (
	"fmt"

	"github.com/jrapoport/gothic/storage"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "migrates the database",
	RunE: migrateRunE,
}

func migrateRunE(*cobra.Command, []string) error {
	c := cmdConfig
	conn, err := storage.Dial(c, c.Log)
	if err != nil {
		err = fmt.Errorf("database error %w", err)
		return err
	}
	if err = conn.MigrateDatabase(); err != nil {
		err = fmt.Errorf("migration error %w", err)
		return err
	}
	c.Log.Infof("migration complete")
	return nil
}
