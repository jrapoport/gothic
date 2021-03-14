package main

/*
import (
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/utils"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:  "delete [email OR id]",
	RunE: deleteUserRunE,
	Args: cobra.ExactArgs(1),
}

func deleteUserRunE(_ *cobra.Command, args []string) error {
	addr, err := mail.ParseAddress(args[0])
	if err != nil {
		return err
	}
	email := addr.Address
	c, err := adminConfig()
	if err != nil {
		return err
	}
	log := c.Log()
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
		log.Warn(err)
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
	yes, err := confirmAction("Delete user %s", email)
	if err != nil {
		return err
	} else if !yes {
		return nil
	}
	err = a.DeleteUser(ctx, u.ID)
	if err != nil {
		err = fmt.Errorf("failed to delete user %s: %w", email, err)
		return err
	}
	fmt.Printf("deleted user: %s\n", email)
	return nil
}
*/
