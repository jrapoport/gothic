package cmd

import (
	"fmt"

	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/models/user"
	"github.com/jrapoport/gothic/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:  "create [email] [password]",
	RunE: createUserRunE,
	Args: cobra.ExactArgs(2),
}

var (
	confirm   bool
	adminRole bool
)

func init() {
	fs := createCmd.Flags()
	fs.BoolVar(&confirm, "confirm", true, "autoconfirm the user")
	fs.BoolVarP(&adminRole, "admin", "a", false, "create admin user (overrides role)")
}

func createUserRunE(_ *cobra.Command, args []string) error {
	email := args[0]
	password := args[1]
	c, err := adminConfig()
	if err != nil {
		return err
	}
	c.Signup.AutoConfirm = confirm
	a, err := core.NewAPI(c)
	if err != nil {
		return err
	}
	ip, err := utils.OutboundIP()
	if err != nil {
		return err
	}
	ctx := context.Background()
	ctx.SetProvider(a.Provider())
	ctx.SetIPAddress(ip.String())
	u, err := a.Signup(ctx, email, "", password, nil)
	if err != nil {
		return err
	}
	if adminRole {
		u, err = a.ChangeRole(ctx, u.ID, user.RoleAdmin)
		if err != nil {
			return err
		}
	}
	fmt.Printf("created user: %s\n", email)
	utils.PrettyPrint(u)
	return nil
}
