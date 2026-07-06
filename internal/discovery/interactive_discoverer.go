package discovery

import (
	"context"
	"net/url"

	"github.com/chromedp/chromedp"
)

type InteractiveDiscoverer struct{}

const jsAlert = "alert(\"You can close this window now\")"

func (InteractiveDiscoverer) discoverMediaManifest(ctx context.Context, pageURL string, urls <-chan *url.URL) (*url.URL, error) {
	if err := chromedp.Run(ctx, chromedp.Navigate(pageURL)); err != nil {
		return nil, err
	}
	manifest, err := waitForURL(ctx, urls)

	if err := chromedp.Run(ctx, chromedp.Evaluate(jsAlert, nil)); err != nil {
		return nil, err
	}

	return manifest, err
}

func (InteractiveDiscoverer) isHeadless() bool {
	return false
}
