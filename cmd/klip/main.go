package main

import (
	"fmt"
	"klip/internal/orchestrator"
	"os"
)

const appName = "Klip"
const errorExitCode = 1

func main() {
	err := orchestrator.Run(appName, os.Args[1:])

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(errorExitCode)
	}
}
