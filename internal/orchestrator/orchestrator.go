package orchestrator

import (
	"context"
	"fmt"
	"klip/internal/cli"
	"klip/internal/core"
	"klip/internal/discovery"
)

func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, args)

	if err != nil {
		return err
	}

	ctx, stop := context.WithTimeout(context.Background(), core.TimeoutValue)
	defer stop()

	media, err := discovery.GetMediaURL(ctx, config.URL)

	if err != nil {
		return err
	}
	fmt.Println(media)
	return nil
}
