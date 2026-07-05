package cli

import (
	"bytes"
	"flag"
	"fmt"
)

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
