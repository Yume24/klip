package discovery

import (
	"context"
	"net/url"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Corresponds to <video>
const videoTag = "video"

// Returns the handler for network events
func captureEventsHandler(ctx context.Context, ch chan<- networkResponse) func(any) {
	return func(event any) {
		go func() {
			if event, ok := event.(*network.EventResponseReceived); ok {
				url, err := url.ParseRequestURI(event.Response.URL)
				if err != nil {
					return
				}

				select {
				case ch <- networkResponse{url: *url, contentType: event.Response.MimeType}:
				case <-ctx.Done():
				}
			}
		}()
	}
}

// Initializes the headless browser and network event capturing
func initializeBrowser(ctx context.Context) (*browserContext, error) {
	eventsChan := make(chan networkResponse)
	ctx, stop := chromedp.NewContext(ctx)

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		stop()
		return nil, err
	}
	chromedp.ListenTarget(ctx, captureEventsHandler(ctx, eventsChan))

	return &browserContext{ctx: ctx, stop: stop, eventsChan: eventsChan}, nil
}

// Clicks on the video tag
func clickVideo(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.Click(videoTag, chromedp.ByQuery))
}

// Navigates to the specified page
func navigateToPage(ctx context.Context, url string) error {
	return chromedp.Run(ctx, chromedp.Navigate(url), chromedp.WaitVisible(videoTag, chromedp.ByQuery))
}
