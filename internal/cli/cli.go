package cli

import (
	"flag"
)

// User supplied config
type config struct {
	URL string
}

func ParseArguments(name string, args []string) (*config, error) {
	config := &config{}
	flagSet := flag.NewFlagSet(name, flag.ExitOnError)
	if err := loadFlagsIntoConfig(config, flagSet, args); err != nil {
		return nil, err
	}
	return config, nil
}
