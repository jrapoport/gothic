package cmd

import (
	"github.com/jrapoport/gothic"
	"github.com/jrapoport/gothic/config"
	"github.com/spf13/cobra"
)

// ExecuteRoot executes the main cmd
func ExecuteRoot() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:     "gothic",
	Version: config.BuildVersion(),
	RunE:    rootRunE,
}

func init() {
	rootCmd.AddCommand(adminCmd)
	rootCmd.AddCommand(serveCmd)
	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c", "", "the config file to use")
}

var configFile = ""

func rootConfig() (*config.Config, error) {
	return config.LoadConfig(configFile)
}

func rootRunE(*cobra.Command, []string) error {
	c, err := rootConfig()
	if err != nil {
		return err
	}
	return gothic.Main(c)
}
