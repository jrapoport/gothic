package user

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "user",
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(roleCmd)
	Cmd.AddCommand(deleteCmd)
}
