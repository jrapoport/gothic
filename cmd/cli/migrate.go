package main

import (
	"fmt"

	"github.com/jrapoport/gothic/store"
	"github.com/spf13/cobra"
)

// this command requires direct DB access
var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "migrates the database",
	RunE: migrateRunE,
}

func init() {
	AddRootCommand(migrateCmd)
}

func migrateRunE(cmd *cobra.Command, _ []string) error {
	c := rootConfig()
	err := c.DB.CheckRequired()
	if err != nil {
		return err
	}
	conn, err := store.Dial(c, nil)
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
