package orchestrator

import (
	"context"
	"fmt"
	"klip/internal/cli"
	"klip/internal/core"
	"klip/internal/discovery"
)

// Entrypoint to the program's flow
func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, args)

	if err != nil {
		return err
	}

	ctx, stop := context.WithTimeout(context.Background(), core.TimeoutValue)
	defer stop()

	media, err := discovery.GetMedia(ctx, config.URL)

	if err != nil {
		return err
	}
	fmt.Println(media.URL.String())
	return nil
}
