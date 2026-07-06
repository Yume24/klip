package cli

import (
	"flag"
	"fmt"
	"klip/internal/core"
)

const urlPosition = 0
const expectedArgumentNumber = 1

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
