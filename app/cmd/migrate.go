package cmd

import (
	"fmt"

	"github.com/jrapoport/gothic/store"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "migrates the database",
	RunE: migrateRunE,
}

func migrateRunE(*cobra.Command, []string) error {
	c, err := adminConfig()
	if err != nil {
		return err
	}
	conn, err := store.Dial(c, c.Log())
	if err != nil {
		return err
	}
	fmt.Println("starting migration...")
	err = conn.AutoMigrate()
	if err != nil {
		return err
	}
	fmt.Println("migration complete")
	return nil
}
