package main

import (
	"fmt"

	"github.com/jrapoport/gothic/cmd/cli/root"
	"github.com/jrapoport/gothic/core/users"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/store"
	"github.com/jrapoport/gothic/utils"
	"github.com/spf13/cobra"
)

// this command requires direct DB access
var passwordCmd = &cobra.Command{
	Use:  "password [OLD PASSWORD] [NEW PASSWORD]",
	Long: "password changes the root password",
	RunE: passwordRunE,
	Args: cobra.ExactArgs(2),
}

func init() {
	root.AddCommand(passwordCmd)
}

func passwordRunE(_ *cobra.Command, args []string) error {
	var (
		oldPassword = args[0]
		newPassword = args[1]
	)
	cfg := root.Config()
	err := cfg.DB.CheckRequired()
	if err != nil {
		return err
	}
	conn, err := store.Dial(cfg, nil)
	if err != nil {
		return err
	}
	fmt.Println("changing root password...")
	su, err := users.GetUser(conn, user.SuperAdminID)
	if err != nil {
		return err
	}
	err = su.Authenticate(oldPassword)
	if err != nil {
		return err
	}
	hash := utils.HashPassword(newPassword)
	err = conn.Model(su).Update("password", hash).Error
	if err != nil {
		return err
	}
	fmt.Println("root password changed")
	return nil
}
