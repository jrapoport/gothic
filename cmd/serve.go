package cmd

import "github.com/spf13/cobra"

var serveCmd = &cobra.Command{
	Use:               "serve",
	Long:              "start server",
	PersistentPreRunE: initConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Main(cmdConfig)
	},
}
