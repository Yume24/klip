package discovery

import (
	"context"
	"fmt"
	"klip/internal/core"
)

// Message returned after the search has timed out
const timeoutErrorMsg = "timeout: could not locate the video after %d seconds"
var timeoutError = fmt.Errorf(timeoutErrorMsg, int(core.TimeoutValue.Seconds()))

func performWebpageFlow(ctx context.Context, pageURL string) error {
	if err := navigateToPage(ctx, pageURL); err != nil {
		return err
	}

	if err := clickVideo(ctx); err != nil {
		return err
	}

	return nil
}

// Function that waits for inspection results
func waitForMedia(ctx context.Context, result <-chan core.Media) (*core.Media, error) {
	select {
	case <-ctx.Done():
		return nil, timeoutError
	case media := <-result:
		return &media, nil
	}
}

// Retruns the Media struct containing information about the discovered media source
func GetMedia(ctx context.Context, pageURL string) (*core.Media, error) {
	browserCtx, err := initializeBrowser(ctx)

	if err != nil {
		return nil, err
	}

	defer browserCtx.stop()
	result := make(chan core.Media)

	go inspectIncomingTraffic(browserCtx.ctx, browserCtx.eventsChan, result)

	if err := performWebpageFlow(browserCtx.ctx, pageURL); err != nil {
		return nil, timeoutError
	}

	return waitForMedia(browserCtx.ctx, result)
}
