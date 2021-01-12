package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string
var Build string

var versionCmd = cobra.Command{
	Run: showVersion,
	Use: "version",
}

func showVersion(*cobra.Command, []string) {
	fmt.Printf("%s (%s)\n", Version, Build)
}
