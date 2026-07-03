package orchestrator

import (
	"context"
	"fmt"
	"klip/internal/cli"
	"klip/internal/discovery"
	"os"
)

func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, os.Args[1:])

	if err != nil {
		return err
	}

	ctx := context.Background()
	media, err := discovery.GetMediaUrl(ctx, config.Url)

	if err != nil {
		return err
	}
	fmt.Println(media)
	return nil
}
