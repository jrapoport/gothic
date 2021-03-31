package main

import (
	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/utils"
	"github.com/spf13/cobra"
)

var (
	configFile string
	cfg        *config.Config
)

var rootCmd = &cobra.Command{
	Use:               utils.ExecutableName(),
	Short:             "control plane for gothic",
	Version:           config.BuildVersion(),
	RunE:              rootRunE,
	PersistentPreRunE: initConfig,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c", "", "the config file to use")
	rootCmd.AddCommand(codeCmd)
	//rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(migrateCmd)
}

func initConfig(cmd *cobra.Command, _ []string) (err error) {
	cfg, err = config.LoadConfig(configFile)
	if err != nil {
		return err
	}
	cfg.DB.AutoMigrate = false
	//cfg.Signup.AutoConfirm = confirm
	cfg.Signup.Default.Username = true
	cfg.Validation.PasswordRegex = ""
	cfg.ReplaceLog(cfg.Log().WithName(cmd.Use))
	return nil
}

func rootRunE(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func rootConfig() *config.Config {
	return cfg
}

// ExecuteRoot executes the main cmd
func ExecuteRoot() error {
	return rootCmd.Execute()
}

/*
func confirmAction(format string, a ...interface{}) (bool, error) {
	p := fmt.Sprintf(format, a...)
	p = fmt.Sprintf("%s? [Yes/No]", p)
	prompt := promptui.Select{
		Label: p,
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return result == "Yes", nil
}

*/
