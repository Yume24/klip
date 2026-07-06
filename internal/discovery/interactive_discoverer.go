package discovery

import (
	"context"

	"github.com/chromedp/chromedp"
)

type InteractiveDiscoverer struct{}

func (InteractiveDiscoverer) discoverMediaManifest(ctx context.Context, pageURL string) error {
	if err := chromedp.Run(ctx, chromedp.Navigate(pageURL)); err != nil {
		return err
	}

	return nil
}
