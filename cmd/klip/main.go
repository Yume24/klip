package main

import (
	"errors"
	"fmt"
	"klip/internal/orchestrator"
	"os"
)

const appName = "Klip"
const errorExitCode = 1
const successExitCode = 0

type exitCoder interface {
	error
	ExitCode() int
}

func main() {
	err := orchestrator.Run(appName, os.Args[1:]) // First arg is the command itself

	if err != nil {
		w := os.Stderr
		code := errorExitCode

		if withCodeError, ok := errors.AsType[exitCoder](err); ok {
			code = withCodeError.ExitCode()
			if code == successExitCode {
				w = os.Stdout
			}
		}
		_, _ = fmt.Fprintln(w, err)
		os.Exit(code)
	}
}
