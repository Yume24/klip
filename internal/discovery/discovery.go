package discovery

import (
	"context"
	"net/url"
)

// Message returned after the search has timed out
const timeoutErrorMsg = "timeout: could not locate the video after %d seconds"

func waitForURL(ctx context.Context, result <-chan *url.URL) (*url.URL, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case url := <-result:
		return url, nil
	}
}

// Returns the URL pointing to video manifest
func DiscoverManifestURL(pageURL string, discoverer Discoverer) (*url.URL, error) {
	browserCtx, cleanup, networkEvents, err := initializeBrowser()
	if err != nil {
		return nil, err
	}

	defer cleanup()

	manifests := make(chan *url.URL)

	go inspectIncomingTraffic(browserCtx, networkEvents, manifests)

	if err := discoverer.discoverMediaManifest(browserCtx, pageURL); err != nil {
		return nil, err
	}

	return waitForURL(browserCtx, manifests)
}
