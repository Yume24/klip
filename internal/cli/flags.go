package cli

import (
	"flag"
	"fmt"
)

const invalidPositionalArgsNumMsgErr = "expected %d arguments, got: %d — type -h or --help for usage"
const positionalArgsNum = 1
const fallbackArgPosition = 0

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
			return fmt.Errorf(invalidPositionalArgsNumMsgErr, positionalArgsNum, flagSet.NArg())
		}
	}

	return nil
}
