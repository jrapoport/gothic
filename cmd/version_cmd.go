package cmd

import (
	"fmt"

	"github.com/jrapoport/gothic/conf"
	"github.com/spf13/cobra"
)

var versionCmd = cobra.Command{
	Run: showVersion,
	Use: "version",
}

func showVersion(*cobra.Command, []string) {
	fmt.Printf("%s\n", conf.CurrentVersion())
}
