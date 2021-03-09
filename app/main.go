package main

import (
	"github.com/jrapoport/gothic/app/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.ExecuteRoot(); err != nil {
		logrus.Fatal(err)
	}
}
