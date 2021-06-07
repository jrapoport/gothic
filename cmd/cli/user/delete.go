package user

import (
	"fmt"
	"net/mail"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/cmd/cli/root"
	"github.com/jrapoport/gothic/core/context"
	"github.com/spf13/cobra"
)

var hard bool

var deleteCmd = &cobra.Command{
	Use:  "delete [ID or EMAIL]",
	RunE: deleteUserRunE,
	Args: cobra.ExactArgs(1),
}

func init() {
	fs := deleteCmd.Flags()
	fs.BoolVarP(&hard, "hard", "h", false, "hard delete user")
}

func deleteUserRunE(_ *cobra.Command, args []string) error {
	client, err := root.NewAdminClient()
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()
	userID := args[0]
	yes := root.ConfirmAction("Delete user %s", userID)
	if !yes {
		return nil
	}
	req := &admin.DeleteUserRequest{
		User: &admin.DeleteUserRequest_UserId{UserId: userID},
		Hard: hard,
	}
	addr, err := mail.ParseAddress(userID)
	if err == nil {
		req.User = &admin.DeleteUserRequest_Email{Email: addr.Address}
	}
	res, err := client.DeleteUser(context.Background(), req)
	if err != nil {
		return err
	}
	fmt.Printf("deleted user: %s\n", res.GetUserId())
	return nil
}
