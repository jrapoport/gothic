package main

import (
	"log"

	"github.com/jrapoport/gothic/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
