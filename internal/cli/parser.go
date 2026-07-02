// Responsible for parsing CLI arguements and returning the Config object
// or an error if any
package cli

import (
	"bytes"
	"flag"
	"fmt"
	"klip/internal/core"
)

// Number of excpeted positional arguments
const argsNum = 1 // We expect only the url to be a positional argument

type Parser struct {
	flagSet    *flag.FlagSet
	args       []string
	usageError *bytes.Buffer
}

func NewParser(name string, osargs []string) *Parser {
	buffer := &bytes.Buffer{}
	flagset := flag.NewFlagSet(name, flag.ContinueOnError)

	flagset.SetOutput(buffer)
	return &Parser{
		flagset,
		osargs,
		buffer,
	}
}

// Parse CLI arguments and return Config object or an error
func (p *Parser) ParseArguments() (*core.Config, error) {
	config := &core.Config{}

	if err := loadFlags(config, p.flagSet, p.args); err != nil {
		return nil, fmt.Errorf("%s", p.usageError) // err is already included in usageError here
	}

	if err := loadArguments(config, argsNum, p.flagSet); err != nil {
		p.flagSet.Usage()
		return nil, fmt.Errorf("%s\n%s", err, p.usageError)
	}

	return config, nil
}
