package orchestrator

import (
	"fmt"
	"klip/internal/cli"
	"klip/internal/strategy"
	"klip/internal/strategy/strategies/hls"
)

func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, args)

	if err != nil {
		return err
	}

	var downloadStrategies []strategy.DownloadStrategy = []strategy.DownloadStrategy{
		&hls.HLSStrategy{},
	}

	downloadStrategy, err := strategy.GetDownloadStrategy(config.URL, downloadStrategies)
	if err != nil {
		return err
	}
	fmt.Println(downloadStrategy)
	return nil
}
