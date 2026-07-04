package discovery

import (
	"context"
	"fmt"
	"klip/internal/core"
)

const timeoutErrorMsg = "timeout: could not locate the video after %d seconds"

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
		return nil, fmt.Errorf(timeoutErrorMsg, int(core.TimeoutValue.Seconds()))
	case <-result:
		return nil, nil
	}
}

func GetMediaURL(ctx context.Context, pageURL string) (*core.Media, error) {
	browserCtx, err := initializeBrowser(ctx)

	if err != nil {
		return nil, err
	}

	defer browserCtx.stop()
	result := make(chan core.Media)

	go inspectIncomingTraffic(browserCtx.ctx, browserCtx.eventsChan, result)

	if err := performWebpageFlow(browserCtx.ctx, pageURL); err != nil {
		return nil, err
	}

	return waitForMedia(browserCtx.ctx, result)
}
