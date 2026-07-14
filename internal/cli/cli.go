package cli

import (
	"flag"
)

// User supplied config
type Config struct {
	URL string
}

func ParseArguments(name string, args []string) (Config, error) {
	config := Config{}
	flagSet := flag.NewFlagSet(name, flag.ExitOnError)
	if err := loadFlagsIntoConfig(&config, flagSet, args); err != nil {
		return Config{}, err
	}
	return config, nil
}
