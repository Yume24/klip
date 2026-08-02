package cli

import (
	"flag"
)

// User supplied config
type Config struct {
	URL      string
	FileName string
}

func ParseArguments(name string, args []string) (config Config, err error) {
	flagSet := flag.NewFlagSet(name, flag.ExitOnError)
	if err = loadFlagsIntoConfig(&config, flagSet, args); err != nil {
		return
	}
	return
}
