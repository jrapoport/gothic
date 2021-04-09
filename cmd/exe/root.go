package main

import (
	"github.com/jrapoport/gothic"
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     utils.ExecutableName(),
	Version: config.BuildVersion(),
	RunE:    rootRunE,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c", "", "the config file to use")
}

var configFile = ""

func initConfig() (*config.Config, error) {
	return config.LoadConfig(configFile)
}

func rootRunE(*cobra.Command, []string) error {
	c, err := initConfig()
	if err != nil {
		return err
	}
	// Call the "real" main in a nested manner so the defers will properly
	// be executed in the case of a graceful shutdown.
	return gothic.Main(c)
}

// ExecuteRoot executes the main cmd
func ExecuteRoot() error {
	return rootCmd.Execute()
}
