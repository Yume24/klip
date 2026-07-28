package cli

import (
	"errors"
	"flag"
	"fmt"
)

const positionalArgsNum = 1
const fallbackArgPosition = 0

var errInvalidPositionalArgsNum = errors.New("invalid number of arguments")

func defineFlags(config *Config, flagSet *flag.FlagSet) {
	flagSet.StringVar(&config.URL, urlFlagName, urlFlagDefaultVal, urlFlagUsage)
	flagSet.StringVar(&config.FileName, fileNameFlagName, fileNameFlagDefaultVal, fileNameFlagUsage)
}

func loadFlagsIntoConfig(config *Config, flagSet *flag.FlagSet, args []string) error {
	defineFlags(config, flagSet)
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	// Fall back to positional arg
	if config.URL == "" {
		if flagSet.NArg() == positionalArgsNum {
			config.URL = flagSet.Arg(fallbackArgPosition)
		} else {
			return fmt.Errorf("%w: excpected %d, got %d", errInvalidPositionalArgsNum, positionalArgsNum, flagSet.NArg())
		}
	}

	return nil
}
