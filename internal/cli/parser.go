package cli

import (
	"bytes"
	"flag"
	"fmt"
)

type parser struct {
	flagSet  *flag.FlagSet
	args     []string
	usageBuf *bytes.Buffer
}

func (p *parser) renderUsage() string {
	p.usageBuf.Reset()
	p.flagSet.Usage()

	return p.usageBuf.String()
}

func (p *parser) wrapErrorWithUsage(e error) error {
	return fmt.Errorf("%w\n%s", e, p.renderUsage())
}

func createParser(name string, osargs []string) *parser {
	buffer := &bytes.Buffer{}
	flags := flag.NewFlagSet(name, flag.ContinueOnError)

	flags.SetOutput(buffer)
	return &parser{
		flagSet:  flags,
		args:     osargs,
		usageBuf: buffer,
	}
}
