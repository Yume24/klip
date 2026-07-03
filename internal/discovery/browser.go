package discovery

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const videoTag = "video"
const contentTypeHeader = "ContentType"

func captureEventsHandler(event any) {
	switch event := event.(type) {
	case *network.EventResponseReceived:
		fmt.Println(event.Response.Headers[contentTypeHeader])
	}
}

func initializeBrowser(ctx context.Context) (context.Context, context.CancelFunc, error) {
	ctx, cancel := chromedp.NewContext(ctx)

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		return nil, cancel, err
	}
	chromedp.ListenTarget(ctx, captureEventsHandler)

	return ctx, cancel, nil
}

func clickVideo(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.Click(videoTag, chromedp.ByQuery))
}

func navigateToPage(ctx context.Context, url string) error {
	return chromedp.Run(ctx, chromedp.Navigate(url))
}
