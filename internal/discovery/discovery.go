package discovery

import (
	"context"
	"errors"
	"fmt"
	"klip/internal/core"
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

func waitForMedia(ctx context.Context, result <-chan core.Media) (*core.Media, error) {
	select {
	case <-ctx.Done():
		return nil, errTimeout
	case media := <-result:
		return &media, nil
	}
}

// Returns the Media struct containing information about the discovered media source
func GetMedia(ctx context.Context, pageURL string) (*core.Media, error) {
	browserCtx, err := initializeBrowser(ctx)

	if err != nil {
		return nil, err
	}

	defer browserCtx.stop()
	result := make(chan core.Media)

	go inspectIncomingTraffic(browserCtx.ctx, browserCtx.eventsChan, result)

	if err := performWebpageFlow(browserCtx.ctx, pageURL); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errTimeout
		}

		return nil, fmt.Errorf("loading page: %w", err)
	}

	return waitForMedia(browserCtx.ctx, result)
}
