package main

import (
	"fmt"

	"github.com/jrapoport/gothic/cmd/cli/root"
	"github.com/jrapoport/gothic/store"
	"github.com/spf13/cobra"
)

// this command requires direct DB access
var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "migrates the database",
	RunE: migrateRunE,
}

func migrateRunE(*cobra.Command, []string) error {
	cfg := root.Config()
	err := cfg.DB.CheckRequired()
	if err != nil {
		return err
	}
	yes := root.ConfirmAction("Start DB migration")
	if !yes {
		return nil
	}
	conn, err := store.Dial(cfg, nil)
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
