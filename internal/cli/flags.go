package cli

import (
	"flag"
	"klip/internal/core"
)

// Defines the CLI flags to be parsed
func defineFlags(flagSet *flag.FlagSet) {

}

// Loads the specified CLI flags into the passed config struct
func loadFlags(config *core.Config, flagSet *flag.FlagSet, args []string) error {
	defineFlags(flagSet)
	return flagSet.Parse(args)
}
