package main

import (
	"fmt"
	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var createCmd = &cobra.Command{
	Use:  "create [email] [password]",
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
	AddRootCommand(createCmd)
}

func createUserRunE(_ *cobra.Command, args []string) error {
	email := args[0]
	password := args[1]
	c := rootConfig()
	c.Signup.AutoConfirm = confirm
	conn, err := clientConn(c.AdminAddress)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
	}()
	if username == randUsername ||
		c.Signup.Username && username == "" ||
		c.Signup.Default.Username && username == "" {
		username = utils.RandomUsername()
	}
	client := admin.NewAdminClient(conn)
	pw := c.RootPassword
	ctx := metadata.NewOutgoingContext(context.Background(),
		metadata.Pairs(rpc.RootPassword, pw))
	res, err := client.CreateUser(ctx, &admin.CreateUserRequest{
		Email:    email,
		Password: password,
		Admin:    adminRole,
		Username: &username,
	})
	if err != nil {
		return err
	}
	fmt.Printf("created %s %s\n",
		res.GetRole(), res.GetEmail())
	return nil
}
