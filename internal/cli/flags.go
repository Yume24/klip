package cli

import (
	"flag"
	"klip/internal/core"
)

const interactiveFlagUsage = "drive the browser interactively to locate the video"
const interactiveFlagDefaultVal = false
const interactiveFlagName = "interactive"

func defineFlags(flagSet *flag.FlagSet, config *core.Config) {
	flagSet.BoolVar(&config.Interactive, interactiveFlagName, interactiveFlagDefaultVal, interactiveFlagUsage)
}

func loadFlagsIntoConfig(config *core.Config, flagSet *flag.FlagSet, args []string) error {
	defineFlags(flagSet, config)
	return flagSet.Parse(args)
}
