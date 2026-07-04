package orchestrator

import (
	"context"
	"fmt"
	"klip/internal/cli"
	"klip/internal/discovery"
	"time"
)

const timeoutValue = time.Second * 60

func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, args)

	if err != nil {
		return err
	}

	ctx, stop := context.WithTimeout(context.Background(), timeoutValue)
	defer stop()

	media, err := discovery.GetMediaUrl(ctx, config.URL)

	if err != nil {
		return err
	}
	fmt.Println(media)
	return nil
}
