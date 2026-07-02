package main

import (
	"fmt"
	"klip/internal/orchestrator"
	"os"
)

const name = "Klip"

func main() {
	err := orchestrator.Run(name, os.Args[1:]) // First arg is the command itself

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
