package orchestrator

import (
	"fmt"
	"klip/internal/cli"
	"klip/internal/discovery"
)

// Entrypoint to the program's flow
func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, args)

	if err != nil {
		return err
	}

	media, err := discovery.GetMedia(config.URL)

	if err != nil {
		return err
	}
	fmt.Println(media.URL.String())
	return nil
}
