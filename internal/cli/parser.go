// Responsible for parsing CLI arguments and returning the Config object
// or an error if any
package cli

import (
	"bytes"
	"flag"
	"fmt"
	"klip/internal/core"
)

// Number of expected positional arguments
const argsNum = 1 // We expect only the url to be a positional argument

type parser struct {
	flagSet    *flag.FlagSet
	args       []string
	usageError *bytes.Buffer
}

func createParser(name string, osargs []string) *parser {
	buffer := &bytes.Buffer{}
	flagset := flag.NewFlagSet(name, flag.ContinueOnError)

	flagset.SetOutput(buffer)
	return &parser{
		flagset,
		osargs,
		buffer,
	}
}

// Parse CLI arguments and return Config object or an error
func ParseArguments(name string, osargs []string) (*core.Config, error) {
	parser := createParser(name, osargs)
	config := &core.Config{}

	if err := loadFlags(config, parser.flagSet, parser.args); err != nil {
		return nil, fmt.Errorf("%s", parser.usageError) // err is already included in usageError here
	}

	if err := loadArguments(config, argsNum, parser.flagSet); err != nil {
		parser.flagSet.Usage()
		return nil, fmt.Errorf("%w\n%s", err, parser.usageError)
	}

	return config, nil
}
