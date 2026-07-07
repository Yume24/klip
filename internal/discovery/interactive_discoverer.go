package discovery

import (
	"context"
	"net/url"

	"github.com/chromedp/chromedp"
)

type InteractiveDiscoverer struct{}

const jsAlert = "alert(\"You can close this window now\")"

func (InteractiveDiscoverer) discoverMediaManifest(ctx context.Context, pageURL string, manifests <-chan *url.URL) (*url.URL, error) {
	if err := chromedp.Run(ctx, chromedp.Navigate(pageURL)); err != nil {
		return nil, err
	}

	var manifest *url.URL
	select {
	case manifest = <-manifests:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if err := chromedp.Run(ctx, chromedp.Evaluate(jsAlert, nil)); err != nil {
		// We can ignore any error at this stage
		return manifest, nil
	}

	<-ctx.Done()
	return manifest, nil
}

func (InteractiveDiscoverer) isHeadless() bool {
	return false
}
