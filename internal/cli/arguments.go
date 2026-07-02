package cli

import (
	"flag"
	"fmt"
	"klip/internal/core"
)

// Load positional arguments into config
// Takes in an excpeted number of positional arguments
func loadArguments(config *core.Config, n int, flagSet *flag.FlagSet) error {
	args, err := getArguments(n, flagSet)

	if err != nil {
		return err
	}

	config.Url = args[0]
	return nil
}

// Gets the remaning positional arguments after parsing
// Takes in an excpeted number of positional arguments
// and returns an error if the actual number does not match.
func getArguments(n int, flagSet *flag.FlagSet) ([]string, error) {
	argumentsLength := flagSet.NArg()

	if argumentsLength != n {
		flagSet.Usage()
		return nil, fmt.Errorf("Excpected %d arguments, got %d", n, argumentsLength)
	}

	return flagSet.Args(), nil
}
