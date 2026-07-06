package discovery

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
)

// Corresponds to <video>
const videoTag = "video"
const timeoutValue = time.Second * 30
const timeoutErrorMsg = "timeout: could not locate the video after %d seconds"

// Message returned after the search has timed out
var errTimeout = fmt.Errorf(timeoutErrorMsg, int(timeoutValue.Seconds()))

type AutoDiscoverer struct{}

func wrapError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return errTimeout
	}

	return fmt.Errorf("loading page: %w", err)
}

// Clicks on the video tag
func clickVideo(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.Click(videoTag, chromedp.ByQuery))
}

// Navigates to the specified page
func navigateToPage(ctx context.Context, url string) error {
	return chromedp.Run(ctx, chromedp.Navigate(url), chromedp.WaitVisible(videoTag, chromedp.ByQuery))
}

func (AutoDiscoverer) discoverMediaManifest(ctx context.Context, pageURL string, urls <-chan *url.URL) (*url.URL, error) {
	timeoutCtx, timeoutCtxStop := context.WithTimeout(ctx, timeoutValue)
	defer timeoutCtxStop()

	if err := navigateToPage(timeoutCtx, pageURL); err != nil {
		return nil, wrapError(err)
	}

	if err := clickVideo(timeoutCtx); err != nil {
		return nil, wrapError(err)
	}

	return waitForURL(timeoutCtx, urls)
}

func (AutoDiscoverer) isHeadless() bool {
	return true
}
