package discovery

import (
	"context"
	"klip/internal/core"
	"mime"
	"net/url"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Corresponds to <video>
const videoTag = "video"
const networkResponseChannelSize = 64

func parseNetworkEvent(event *network.EventResponseReceived) (*networkResponse, error) {
	requestURL, err := url.ParseRequestURI(event.Response.URL)
	if err != nil {
		return nil, err
	}

	contentType, _, _ := mime.ParseMediaType(event.Response.MimeType)

	return &networkResponse{url: requestURL, contentType: contentType}, nil
}

func sendResponseToChannel(ctx context.Context, ch chan<- networkResponse, resp *networkResponse) {
	select {
	case ch <- *resp:
	case <-ctx.Done():
	}
}

// Returns the handler for network events
func captureEventsHandler(ctx context.Context, ch chan<- networkResponse) func(any) {
	return func(event any) {
		if event, ok := event.(*network.EventResponseReceived); ok {
			resp, err := parseNetworkEvent(event)

			if err != nil {
				return
			}

			go sendResponseToChannel(ctx, ch, resp)
		}
	}
}

func initializeContext() (context.Context, context.CancelFunc) {
	timeoutCtx, stopTimeoutCtx := context.WithTimeout(context.Background(), core.TimeoutValue)
	ctx, stopBrowserCtx := chromedp.NewContext(timeoutCtx)

	cleanup := func() {
		stopBrowserCtx()
		stopTimeoutCtx()
	}

	return ctx, cleanup
}

// Initializes the headless browser and network event capturing
func initializeBrowser() (context.Context, context.CancelFunc, <-chan networkResponse, error) {
	eventsChan := make(chan networkResponse, networkResponseChannelSize)

	ctx, cleanup := initializeContext()

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		cleanup()
		return nil, nil, nil, err
	}
	chromedp.ListenTarget(ctx, captureEventsHandler(ctx, eventsChan))

	return ctx, cleanup, eventsChan, nil
}

// Clicks on the video tag
func clickVideo(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.Click(videoTag, chromedp.ByQuery))
}

// Navigates to the specified page
func navigateToPage(ctx context.Context, url string) error {
	return chromedp.Run(ctx, chromedp.Navigate(url), chromedp.WaitVisible(videoTag, chromedp.ByQuery))
}
