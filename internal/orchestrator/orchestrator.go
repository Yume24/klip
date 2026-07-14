package orchestrator

import (
	"github.com/Yume24/klip/internal/cli"
	"github.com/Yume24/klip/internal/strategy"
	"github.com/Yume24/klip/internal/strategy/strategies/hls"
)

func createDownloadStrategies() []strategy.DownloadStrategy {
	return []strategy.DownloadStrategy{
		&hls.HLSStrategy{},
	}
}

func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, args)

	if err != nil {
		return err
	}

	downloadStrategies := createDownloadStrategies()

	downloadStrategy, err := strategy.GetDownloadStrategy(config.URL, downloadStrategies)
	if err != nil {
		return err
	}

	err = downloadStrategy.Download()
	if err != nil {
		return err
	}
	return nil
}
