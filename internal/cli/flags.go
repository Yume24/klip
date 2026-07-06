package cli

import (
	"flag"
	"klip/internal/core"
)

func defineFlags(flagSet *flag.FlagSet) {

}

func loadFlagsIntoConfig(config *core.Config, flagSet *flag.FlagSet, args []string) error {
	defineFlags(flagSet)
	return flagSet.Parse(args)
}
