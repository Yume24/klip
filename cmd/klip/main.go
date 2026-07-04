package main

import (
	"errors"
	"flag"
	"fmt"
	"klip/internal/orchestrator"
	"os"
)

const appName = "Klip"
const successExitCode = 0
const errorExitCode = 1

func main() {
	err := orchestrator.Run(appName, os.Args[1:]) // First arg is the command itself

	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		if errors.Is(err, flag.ErrHelp) {
			os.Exit(successExitCode)
		}

		os.Exit(errorExitCode)
	}
}
