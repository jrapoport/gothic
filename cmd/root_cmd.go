package cmd

import (
	"github.com/jrapoport/gothic/conf"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configFile = ""

var rootCmd = cobra.Command{
	Use: "gothic",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

// RootCommand will setup and return the root command
func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&serveCmd, &migrateCmd, &versionCmd, adminCmd())
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "the config file to use")

	return &rootCmd
}

func execWithConfig(cmd *cobra.Command, fn func(globalConfig *conf.Configuration)) {
	globalConfig, err := conf.LoadConfiguration(configFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %+v", err)
	}
	fn(globalConfig)
}

func execWithConfigAndArgs(cmd *cobra.Command, fn func(globalConfig *conf.Configuration, args []string), args []string) {
	globalConfig, err := conf.LoadConfiguration(configFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %+v", err)
	}
	fn(globalConfig, args)
}
