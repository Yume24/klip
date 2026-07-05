// Responsible for parsing CLI arguments and returning the Config object
// or an error if any
package cli

import (
	"errors"
	"flag"
	"klip/internal/core"
)

// Number of expected positional arguments
const argsNum = 1 // We expect only the url to be a positional argument


// Parse CLI arguments and return Config object or an error
func ParseArguments(name string, osargs []string) (*core.Config, error) {
	parser := createParser(name, osargs)
	config := &core.Config{}

	if err := loadFlagsIntoConfig(config, parser.flagSet, parser.args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, &helpError{usageMessage: parser.usage()}
		}
		return nil, parser.wrapErrorWithUsage(err)
	}

	if err := loadArgumentsIntoConfig(config, argsNum, parser.flagSet); err != nil {
		return nil, parser.wrapErrorWithUsage(err)
	}

	return config, nil
}
