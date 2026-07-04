package discovery

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const videoTag = "video"

func captureEventsHandler(ctx context.Context, ch chan<- networkResponse) func(any) {
	return func(event any) {
		go func() {
			if event, ok := event.(*network.EventResponseReceived); ok {
				select {
				case ch <- networkResponse{url: event.Response.URL, contentType: event.Response.MimeType}:
				case <-ctx.Done():
				}
			}
		}()
	}
}

func initializeBrowser(ctx context.Context) (*browserContext, error) {
	eventsChan := make(chan networkResponse)
	ctx, stop := chromedp.NewContext(ctx)

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		stop()
		return nil, err
	}
	chromedp.ListenTarget(ctx, captureEventsHandler(ctx, eventsChan))

	return &browserContext{ctx, stop, eventsChan}, nil
}

func clickVideo(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.Click(videoTag, chromedp.ByQuery))
}

func navigateToPage(ctx context.Context, url string) error {
	return chromedp.Run(ctx, chromedp.Navigate(url), chromedp.WaitVisible(videoTag, chromedp.ByQuery))
}
