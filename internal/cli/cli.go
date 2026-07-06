package cli

import (
	"errors"
	"flag"
	"klip/internal/core"
)

func ParseArguments(name string, osargs []string) (*core.Config, error) {
	parser := createParser(name, osargs)
	config := &core.Config{}

	if err := loadFlagsIntoConfig(config, parser.flagSet, parser.args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, &helpError{usageMessage: parser.renderUsage()}
		}
		return nil, parser.wrapErrorWithUsage(err)
	}

	if err := loadURLIntoConfig(config, parser.flagSet); err != nil {
		return nil, parser.wrapErrorWithUsage(err)
	}

	return config, nil
}
