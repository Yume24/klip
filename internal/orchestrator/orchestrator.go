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

	url, err := discovery.DiscoverManifestURL(config.URL)

	if err != nil {
		return err
	}
	fmt.Println(url.String())
	return nil
}
