package main

import (
	"fmt"
	"klip/internal/orchestrator"
	"os"
)

// App name
const name = "Klip"

func main() {
	err := orchestrator.Run(name, os.Args[1:]) // First arg is the command itself

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
