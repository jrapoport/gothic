package cmd

import "github.com/spf13/cobra"

var adminCmd = &cobra.Command{
	Use:               "admin",
	PersistentPreRunE: initConfig,
}

func init() {
	adminCmd.AddCommand(codeCmd)
	adminCmd.AddCommand(userCmd)
	adminCmd.AddCommand(migrateCmd)
}
