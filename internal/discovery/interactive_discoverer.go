package discovery

import (
	"context"
	"net/url"

	"github.com/chromedp/chromedp"
)

type InteractiveDiscoverer struct{}

const jsAlert = "alert(\"Done! Click OK to close\")"

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

	// alert() blocks until the user dismisses it
	// that dismissal is the signal to close the browser and finish
	_ = chromedp.Run(ctx, chromedp.Evaluate(jsAlert, nil))
	return manifest, nil
}

func (InteractiveDiscoverer) isHeadless() bool {
	return false
}
