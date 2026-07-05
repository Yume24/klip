package cli

import (
	"flag"
	"fmt"
	"klip/internal/core"
)

// Position of the URL param in positional args
const urlPosition = 0
const expectedArgumentNumber = 1

// Load positional arguments into config
// Takes in an expected number of positional arguments
func loadURLIntoConfig(config *core.Config, flagSet *flag.FlagSet) error {

	if err := validateArgumentsLength(flagSet); err != nil {
		return err
	}

	config.URL = flagSet.Arg(urlPosition)
	return nil
}

func validateArgumentsLength(flagSet *flag.FlagSet) error {
	argumentsLength := flagSet.NArg()

	if argumentsLength != expectedArgumentNumber {
		return fmt.Errorf("expected %d argument(s), got %d", expectedArgumentNumber, argumentsLength)
	}

	return nil
}
