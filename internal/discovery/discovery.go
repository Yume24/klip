package discovery

import (
	"context"
	"net/url"
)

// Returns the URL pointing to video manifest
func DiscoverManifestURL(pageURL string, discoverer Discoverer) (*url.URL, error) {
	browserCtx, cleanup, networkEvents, err := initializeBrowser(discoverer.isHeadless())
	if err != nil {
		return nil, err
	}

	defer cleanup()

	manifests := make(chan *url.URL)

	go inspectIncomingTraffic(browserCtx, networkEvents, manifests)

	return discoverer.discoverMediaManifest(browserCtx, pageURL, manifests)
}

func waitForURL(ctx context.Context, result <-chan *url.URL) (*url.URL, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case url := <-result:
		return url, nil
	}
}
