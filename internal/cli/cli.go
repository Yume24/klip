package cli

import (
	"flag"
	"klip/internal/core"
)

func ParseArguments(name string, args []string) (*core.Config, error) {
	config := &core.Config{}
	flagSet := flag.NewFlagSet(name, flag.ExitOnError)
	if err := loadFlagsIntoConfig(config, flagSet, args); err != nil {
		return nil, err
	}
	return config, nil
}
