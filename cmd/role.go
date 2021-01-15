package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:  "role [email or id]",
	RunE: roleUserRunE,
	Args: cobra.ExactArgs(1),
}

func init() {
	fs := roleCmd.Flags()
	fs.StringVarP(&userRole, "role", "r", "user", "set the role for new user")
	fs.BoolVarP(&superRole, "super-admin", "s", false, "create super admin user (overrides admin & role)")
	fs.BoolVarP(&adminRole, "admin", "a", false, "create admin user (overrides role)")
}

func roleUserRunE(_ *cobra.Command, args []string) error {
	email := args[0]
	c := cmdConfig
	db, err := storage.Dial(c, c.Log)
	if err != nil {
		return err
	}
	user, err := models.FindUserByEmail(db, email)
	if err != nil {
		c.Log.Warn(err)
		var uid uuid.UUID
		uid, err = uuid.Parse(email)
		if err != nil {
			return err
		}
		user, err = models.FindUserByID(db, uid)
	}
	if err != nil {
		return err
	} else if user == nil {
		err = fmt.Errorf("user not found: %s", role(c))
		return err
	}
	user.Role = role(c)
	user.IsSuperAdmin = superRole
	err = db.Model(&user).Select("role", "is_super_admin").Updates(user).Error
	if err != nil {
		err = fmt.Errorf("failed to update user role %s: %w", email, err)
		return err
	}
	c.Log.Infof("updated user role %s: %s (super admin: %t)",
		email, user.Role, user.IsSuperAdmin)
	return nil
}
