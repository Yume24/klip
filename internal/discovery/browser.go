package discovery

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const headlessFlag = "headless"

type networkEventHandler func(*network.EventResponseReceived)

func createNetworkEventsHandler(handler networkEventHandler) func(any) {
	return func(event any) {
		if event, ok := event.(*network.EventResponseReceived); ok {
			handler(event)
		}
	}
}

func initializeContext(isHeadless bool) (context.Context, context.CancelFunc) {
	options := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag(headlessFlag, isHeadless))
	allocCtx, stopAllocCtx := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, stopCtx := chromedp.NewContext(allocCtx)

	cleanup := func() {
		stopCtx()
		stopAllocCtx()
	}

	return ctx, cleanup
}

// Initializes the headless browser and network event capturing
func initializeBrowser(isHeadless bool, networkEventHandler networkEventHandler) (context.Context, context.CancelFunc, error) {
	ctx, cleanup := initializeContext(isHeadless)

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		cleanup()
		return nil, nil, err
	}
	chromedp.ListenTarget(ctx, createNetworkEventsHandler(networkEventHandler))

	return ctx, cleanup, nil
}
