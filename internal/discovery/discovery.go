package discovery

import (
	"context"
	"errors"
	"fmt"
	"klip/internal/core"
	"net/url"
)

// Message returned after the search has timed out
const timeoutErrorMsg = "timeout: could not locate the video after %d seconds"

var errTimeout = fmt.Errorf(timeoutErrorMsg, int(core.TimeoutValue.Seconds()))

func performWebpageFlow(ctx context.Context, pageURL string) error {
	if err := navigateToPage(ctx, pageURL); err != nil {
		return err
	}

	if err := clickVideo(ctx); err != nil {
		return err
	}

	return nil
}

func waitForURL(ctx context.Context, result <-chan *url.URL) (*url.URL, error) {
	select {
	case <-ctx.Done():
		return nil, errTimeout
	case url := <-result:
		return url, nil
	}
}

// Returns the URL pointing to video manifest
func DiscoverManifestURL(pageURL string) (*url.URL, error) {
	browserCtx, cleanup, networkEvents, err := initializeBrowser()
	if err != nil {
		return nil, err
	}

	defer cleanup()

	manifests := make(chan *url.URL)

	go inspectIncomingTraffic(browserCtx, networkEvents, manifests)

	if err := performWebpageFlow(browserCtx, pageURL); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errTimeout
		}

		return nil, fmt.Errorf("loading page: %w", err)
	}

	return waitForURL(browserCtx, manifests)
}
