package orchestrator

import (
	"fmt"
	"klip/internal/cli"
	"klip/internal/strategy"
)

func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, args)

	if err != nil {
		return err
	}

	downloadStrategy, err := strategy.GetDownloadStrategy(config.URL)
	if err != nil {
		return err
	}
	fmt.Println(downloadStrategy)
	return nil
}
