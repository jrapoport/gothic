package cmd

import (
	"fmt"

	"github.com/jrapoport/gothic/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use: "admin",
}

func init() {
	adminCmd.AddCommand(codeCmd)
	adminCmd.AddCommand(userCmd)
	adminCmd.AddCommand(migrateCmd)
}

func adminConfig() (*config.Config, error) {
	c, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, err
	}
	c.DB.AutoMigrate = false
	c.Signup.AutoConfirm = confirm
	c.Signup.Default.Username = true
	c.Validation.PasswordRegex = ""
	return c, nil
}

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
