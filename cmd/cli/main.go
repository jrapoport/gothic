package main

import (
	"fmt"
)

func main() {
	if err := ExecuteRoot(); err != nil {
		fmt.Printf("Error: %s\n\n", err)
	}
}
