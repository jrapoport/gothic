package main

import (
	"fmt"

	"github.com/jrapoport/gothic/cmd/cli/root"
	"github.com/jrapoport/gothic/cmd/cli/user"
)

func init() {
	root.AddCommand(user.Cmd)
	root.AddCommand(codeCmd)
	root.AddCommand(migrateCmd)
}

func main() {
	if err := root.Execute(); err != nil {
		fmt.Printf("Error: %s\n\n", err)
	}
}
