package cli

import (
	"flag"
	"fmt"
	"klip/internal/core"
)

const urlPosition = 0

// Load positional arguments into config
// Takes in an expected number of positional arguments
func loadArguments(config *core.Config, n int, flagSet *flag.FlagSet) error {
	args, err := getArguments(n, flagSet)

	if err != nil {
		return err
	}

	config.Url = args[urlPosition]
	return nil
}

// Gets the remaining positional arguments after parsing
// Takes in an expected number of positional arguments
// and returns an error if the actual number does not match.
func getArguments(n int, flagSet *flag.FlagSet) ([]string, error) {
	argumentsLength := flagSet.NArg()

	if argumentsLength != n {
		return nil, fmt.Errorf("expected %d argument(s), got %d", n, argumentsLength)
	}

	return flagSet.Args(), nil
}
