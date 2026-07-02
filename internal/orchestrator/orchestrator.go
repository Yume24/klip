package orchestrator

import (
	"fmt"
	"klip/internal/cli"
	"os"
)

func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, os.Args[1:])

	if err != nil {
		return err
	}

	fmt.Println(config)
	return nil
}
