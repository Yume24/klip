package discovery

import (
	"context"
	"klip/internal/core"
)

func GetMediaUrl(ctx context.Context, pageUrl string) (*core.Media, error) {
	browserContext, err := initializeBrowser(ctx)

	if err != nil {
		return nil, err
	}

	defer browserContext.cancel()

	go inspectIncomingTraffic(browserContext.ctx, browserContext.eventsChan, browserContext.cancel)

	if err := navigateToPage(browserContext.ctx, pageUrl); err != nil {
		return nil, err
	}

	if err := clickVideo(browserContext.ctx); err != nil {
		return nil, err
	}

	return nil, nil
}
