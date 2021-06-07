package user

import (
	"fmt"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/cmd/cli/root"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:  "create [EMAIL] [PASSWORD]",
	RunE: createUserRunE,
	Args: cobra.ExactArgs(2),
}

const randUsername = "random"

var (
	confirm   bool
	adminRole bool
	username  string
)

func init() {
	fs := createCmd.Flags()
	fs.StringVar(&username, "username", randUsername, "username for the user")
	fs.BoolVar(&confirm, "confirm", true, "autoconfirm the user")
	fs.BoolVarP(&adminRole, "admin", "a", false, "create admin user")
}

func createUserRunE(_ *cobra.Command, args []string) error {
	email := args[0]
	password := args[1]
	cfg := root.Config()
	cfg.Signup.AutoConfirm = confirm
	client, err := root.NewAdminClient()
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()
	if username == randUsername ||
		cfg.Signup.Username && username == "" ||
		cfg.Signup.Default.Username && username == "" {
		username = utils.RandomUsername()
	}
	req := &admin.CreateUserRequest{
		Email:    email,
		Password: password,
		Admin:    adminRole,
		Username: &username,
	}
	res, err := client.CreateUser(context.Background(), req)
	if err != nil {
		return err
	}
	fmt.Printf("created %s: %s (%s)\n",
		res.GetRole(), res.GetUserId(), res.GetEmail())
	return nil
}
