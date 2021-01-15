package cmd

import (
	"errors"
	"fmt"

	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:  "create [email] [password]",
	RunE: createUserRunE,
	Args: cobra.ExactArgs(2),
}

var (
	confirm   bool
	superRole bool
	adminRole bool
	userRole  string
)

func init() {
	fs := createCmd.Flags()
	fs.BoolVar(&confirm, "confirm", false, "confirm the new user")
	fs.StringVarP(&userRole, "role", "r", "user", "set the role for new user")
	fs.BoolVarP(&superRole, "super-admin", "s", false, "create super admin user (overrides admin & role)")
	fs.BoolVarP(&adminRole, "admin", "a", false, "create admin user (overrides role)")
}

func role(c *conf.Configuration) string {
	if superRole || adminRole {
		return c.JWT.AdminGroup
	} else if userRole != "" {
		return userRole
	}
	return c.JWT.DefaultGroup
}

func createUserRunE(_ *cobra.Command, args []string) error {
	email := args[0]
	password := args[1]
	c := cmdConfig
	db, err := storage.Dial(c, c.Log)
	if err != nil {
		return err
	}
	exists, err := models.IsDuplicatedEmail(db, email, "")
	if err != nil {
		return err
	} else if exists {
		return errors.New("user already exists")
	}
	user, err := models.NewUser(email, password, "", nil)
	if err != nil {
		return err
	}
	err = db.Transaction(func(tx *storage.Connection) error {
		user.IsSuperAdmin = superRole
		if err = tx.Create(user).Error; err != nil {
			return err
		}
		if err = user.SetRole(tx, role(c)); err != nil {
			return err
		}
		if !c.Mailer.Autoconfirm && !confirm {
			return nil
		}
		return user.Confirm(tx)
	})
	if err != nil {
		err = fmt.Errorf("failed to create user %s: %w", email, err)
		return err
	}
	c.Log.Infof("created user: %s", email)
	return nil
}
