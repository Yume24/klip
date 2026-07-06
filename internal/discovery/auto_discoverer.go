package discovery

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

// Corresponds to <video>
const videoTag = "video"
const timeoutValue = time.Second * 30

var errTimeout = fmt.Errorf(timeoutErrorMsg, int(timeoutValue.Seconds()))

type AutoDiscoverer struct{}

func wrapError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return errTimeout
	}

	return fmt.Errorf("loading page: %w", err)
}

func (AutoDiscoverer) discoverMediaManifest(ctx context.Context, pageURL string) error {
	timeoutCtx, timeoutCtxStop := context.WithTimeout(ctx, timeoutValue)
	defer timeoutCtxStop()

	if err := navigateToPage(timeoutCtx, pageURL); err != nil {
		return wrapError(err)
	}

	if err := clickVideo(timeoutCtx); err != nil {
		return wrapError(err)
	}

	return nil
}

// Clicks on the video tag
func clickVideo(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.Click(videoTag, chromedp.ByQuery))
}

// Navigates to the specified page
func navigateToPage(ctx context.Context, url string) error {
	return chromedp.Run(ctx, chromedp.Navigate(url), chromedp.WaitVisible(videoTag, chromedp.ByQuery))
}
