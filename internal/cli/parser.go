// Responsible for parsing CLI arguements and returning the Config object
// or an error if any
package cli

import (
	"flag"
	"klip/internal/core"
)

// Number of excpeted positional arguments
const argsNum = 1 // We expect only the url to be a positional argument

type Parser struct {
	FlagSet *flag.FlagSet
	Args    []string
}

func NewParser(name string, osargs []string) *Parser {
	return &Parser{
		flag.NewFlagSet(name, flag.ContinueOnError),
		osargs,
	}
}

// Parse CLI arguments and return Config object or an error
func (p *Parser) ParseArguments() (*core.Config, error) {
	config := &core.Config{}

	if err := loadFlags(config, p.FlagSet, p.Args); err != nil {
		return nil, err
	}

	if err := loadArguments(config, argsNum, p.FlagSet); err != nil {
		return nil, err
	}

	return config, nil
}
