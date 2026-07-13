package cli

import (
	"errors"
	"flag"
	"fmt"
)

const invalidPositionalArgsNumMsg = "invalid number of arguments"
const positionalArgsNum = 1
const fallbackArgPosition = 0

var invalidPositionalArgsNumError = errors.New(invalidPositionalArgsNumMsg)

func defineFlags(config *config, flagSet *flag.FlagSet) {
	flagSet.StringVar(&config.URL, urlFlagName, urlFlagDefaultVal, urlFlagUsage)
}

func loadFlagsIntoConfig(config *config, flagSet *flag.FlagSet, args []string) error {
	defineFlags(config, flagSet)
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	// Fall back to positional arg
	if config.URL == "" {
		if flagSet.NArg() == positionalArgsNum {
			config.URL = flagSet.Arg(fallbackArgPosition)
		} else {
			return fmt.Errorf("%w: excpected %d, got %d", invalidPositionalArgsNumError, positionalArgsNum, flagSet.NArg())
		}
	}

	return nil
}
