package discovery

import (
	"context"
	"klip/internal/core"
)

func GetMediaURL(ctx context.Context, pageURL string) (*core.Media, error) {
	browserCtx, err := initializeBrowser(ctx)

	if err != nil {
		return nil, err
	}

	defer browserCtx.stop()

	go inspectIncomingTraffic(browserCtx.ctx, browserCtx.eventsChan)

	if err := navigateToPage(browserCtx.ctx, pageURL); err != nil {
		return nil, err
	}

	if err := clickVideo(browserCtx.ctx); err != nil {
		return nil, err
	}

	return nil, nil
}
