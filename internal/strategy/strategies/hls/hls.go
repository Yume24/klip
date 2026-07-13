package hls

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const manifestsChanSize = 1

type HLSStrategy struct {
	url string
}

func createNetworkEventHandler(ch chan<- string) func(any) {
	return func(event any) {
		if event, ok := event.(*network.EventResponseReceived); ok {
			if isMediaManifest(event) {
				manifestURL := event.Response.URL
				select {
				case ch <- manifestURL:
				default:
				}
			}
		}
	}
}

func (s *HLSStrategy) Scout(ctx context.Context, pageURL string) bool {
	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		return false
	}

	manifests := make(chan string, manifestsChanSize)

	chromedp.ListenTarget(ctx, createNetworkEventHandler(manifests))

	if err := chromedp.Run(ctx, chromedp.Navigate(pageURL)); err != nil {
		return false
	}

	select {
	case manifest := <-manifests:
		s.url = manifest
		return true
	case <-ctx.Done():
		return false
	}
}

func (s *HLSStrategy) Download() {}
