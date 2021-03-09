package cmd

import "github.com/spf13/cobra"

var serveCmd = &cobra.Command{
	Use:  "serve",
	Long: "start server",
	RunE: rootRunE,
}
