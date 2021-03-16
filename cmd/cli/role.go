package main

/*
import (
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/utils"
	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:  "role [email or id]",
	RunE: roleUserRunE,
	Args: cobra.ExactArgs(1),
}

var userRole string

func init() {
	fs := roleCmd.Flags()
	fs.StringVarP(&userRole, "role", "r", "user", "set the role for new user")
}

func roleUserRunE(_ *cobra.Command, args []string) error {
	addr, err := mail.ParseAddress(args[0])
	if err != nil {
		return err
	}
	email := addr.Address
	c, err := adminConfig()
	if err != nil {
		return err
	}
	a, err := core.NewAPI(c)
	if err != nil {
		return err
	}
	ip, err := utils.OutboundIP()
	if err != nil {
		return err
	}
	ctx := context.Background()
	ctx.SetIPAddress(ip.String())
	u, err := a.GetUserWithEmail(email)
	if err != nil {
		c.Log().Warn(err)
		var uid uuid.UUID
		uid, err = uuid.Parse(email)
		if err != nil {
			return err
		}
		u, err = a.GetUser(uid)
	}
	if err != nil {
		return err
	}
	r := user.RoleUser
	if userRole == user.RoleAdmin.String() {
		r = user.RoleAdmin
	}
	if u.Role == r {
		fmt.Printf("user: %s\nrole: %s\n",
			u.ID, u.Role.String())
		return nil
	}
	yes, err := confirmAction("Change user role to %s [%s]",
		r.String(), u.ID)
	if err != nil {
		return err
	} else if !yes {
		return nil
	}
	_, err = a.ChangeRole(ctx, u.ID, r)
	if err != nil {
		err = fmt.Errorf("failed to delete user %s: %w", email, err)
		return err
	}
	fmt.Printf("user: %s\nrole: %s\n",
		u.ID, u.Role.String())
	return nil
}


*/
