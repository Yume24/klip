package strategy

import (
	"klip/internal/strategy/browser"

	"github.com/chromedp/chromedp"
)

func smokeTest(b *browser.Browser, pageURL string, errors chan<- error) {
	ctx, cancel := b.NewTab(timeout)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(pageURL)); err != nil {
		errors <- err
	}
}
