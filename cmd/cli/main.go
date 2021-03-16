package main

import "github.com/sirupsen/logrus"

func main() {
	if err := ExecuteRoot(); err != nil {
		logrus.Fatal(err)
	}
}
