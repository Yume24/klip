package discovery

import (
	"context"
	"klip/internal/core"
)

func GetMediaUrl(ctx context.Context, pageUrl string) (*core.Media, error) {
	browserCtx, err := initializeBrowser(ctx)

	if err != nil {
		return nil, err
	}

	defer browserCtx.stop()

	go inspectIncomingTraffic(browserCtx.ctx, browserCtx.eventsChan)

	if err := navigateToPage(browserCtx.ctx, pageUrl); err != nil {
		return nil, err
	}

	if err := clickVideo(browserCtx.ctx); err != nil {
		return nil, err
	}

	return nil, nil
}
