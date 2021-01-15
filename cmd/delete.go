package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:  "delete [email OR id]",
	RunE: deleteUserRunE,
	Args: cobra.ExactArgs(1),
}

func deleteUserRunE(_ *cobra.Command, args []string) error {
	email := args[0]
	c := cmdConfig
	db, err := storage.Dial(c, c.Log)
	if err != nil {
		return err
	}
	user, err := models.FindUserByEmailAndAudience(db, email, "")
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
		err = fmt.Errorf("user not found: %s", email)
		return err
	}
	err = db.Delete(user).Error
	if err != nil {
		err = fmt.Errorf("failed to delete user %s: %w", email, err)
		return err
	}
	c.Log.Infof("deleted user: %s", email)
	return nil
}
