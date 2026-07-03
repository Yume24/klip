package discovery

import (
	"context"
	"klip/internal/core"
)

func GetMediaUrl(ctx context.Context, pageUrl string) (*core.Media, error) {
	ctx, cancel, err := initializeBrowser(ctx)
	defer cancel()

	if err != nil {
		return nil, err
	}

	if err := navigateToPage(ctx, pageUrl); err != nil {
		return nil, err
	}

	if err := clickVideo(ctx); err != nil {
		return nil, err
	}

	return nil, nil
}
