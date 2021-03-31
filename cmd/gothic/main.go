package main

import "github.com/jrapoport/gothic/log"

func main() {
	if err := ExecuteRoot(); err != nil {
		log.Fatal(err)
	}
}
