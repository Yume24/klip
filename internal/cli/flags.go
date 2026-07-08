package cli

import (
	"errors"
	"flag"
	"klip/internal/core"
)

const noURLMessageError = "missing URL: Type -h or --help for usage"
const fallbackArgPosition = 0

func defineFlags(config *core.Config, flagSet *flag.FlagSet) {
	flagSet.BoolVar(&config.Interactive, interactiveFlagName, interactiveFlagDefaultVal, interactiveFlagUsage)
	flagSet.StringVar(&config.URL, urlFlagName, urlFlagDefaultVal, urlFlagUsage)
}

func loadFlagsIntoConfig(config *core.Config, flagSet *flag.FlagSet, args []string) error {
	defineFlags(config, flagSet)
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	// Fall back to positional arg
	if config.URL == "" {
		if flagSet.NArg() > 0 {
			config.URL = flagSet.Arg(fallbackArgPosition)
		} else {
			return errors.New(noURLMessageError)
		}
	}

	return nil
}
