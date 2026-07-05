// Responsible for parsing CLI arguments and returning the Config object
// or an error if any
package cli

import (
	"bytes"
	"errors"
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

func (p *parser) usage() string {
	p.usageError.Reset()
	p.flagSet.Usage()

	return p.usageError.String()
}

func (p *parser) wrapErrorWithUsage(e error) error {
	return fmt.Errorf("%w\n%s", e, p.usage())
}

func createParser(name string, osargs []string) *parser {
	buffer := &bytes.Buffer{}
	flags := flag.NewFlagSet(name, flag.ContinueOnError)

	flags.SetOutput(buffer)
	return &parser{
		flagSet:    flags,
		args:       osargs,
		usageError: buffer,
	}
}

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
