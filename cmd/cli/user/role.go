package user

import (
	"fmt"
	"net/mail"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/cmd/cli/root"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/user"
	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:  `role [ID or EMAIL] [ROLE]`,
	RunE: roleUserRunE,
	Args: cobra.ExactArgs(2),
}

var validRoles = []string{
	user.RoleUser.String(),
	user.RoleAdmin.String(),
}

func roleUserRunE(_ *cobra.Command, args []string) error {
	client, err := root.NewAdminClient()
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()
	userID := args[0]
	role := user.ToRole(args[1])
	var valid bool
	for _, validRole := range validRoles {
		if valid = validRole == role.String(); valid {
			break
		}
	}
	if !valid {
		return fmt.Errorf(`'%s' is not a valid role: %v`,
			args[1], validRoles)
	}
	yes := root.ConfirmAction("Change user %s role to %s",
		userID, role.String())
	if !yes {
		return nil
	}
	req := &admin.ChangeUserRoleRequest{
		User: &admin.ChangeUserRoleRequest_UserId{UserId: userID},
		Role: role.String(),
	}
	addr, err := mail.ParseAddress(userID)
	if err == nil {
		req.User = &admin.ChangeUserRoleRequest_Email{Email: addr.Address}
	}
	res, err := client.ChangeUserRole(context.Background(), req)
	if err != nil {
		return err
	}
	fmt.Printf("changed user %s role to: %s\n",
		res.GetUserId(), res.GetRole())
	return nil
}
