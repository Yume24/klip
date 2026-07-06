package orchestrator

import (
	"fmt"
	"klip/internal/cli"
	"klip/internal/core"
	"klip/internal/discovery"
)

func decideDiscoveryStrategy(config *core.Config) discovery.Discoverer {
	if config.Interactive {
		return discovery.InteractiveDiscoverer{}
	}

	return discovery.AutoDiscoverer{}
}

// Entrypoint to the program's flow
func Run(name string, args []string) error {
	config, err := cli.ParseArguments(name, args)

	if err != nil {
		return err
	}

	discoverer := decideDiscoveryStrategy(config)
	url, err := discovery.DiscoverManifestURL(config.URL, discoverer)

	if err != nil {
		return err
	}
	fmt.Println(url.String())
	return nil
}
