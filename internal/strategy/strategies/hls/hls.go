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

func (s *HLSStrategy) Scout(ctx context.Context, pageURL string) (bool, error) {
	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		return false, err
	}

	manifests := make(chan string, manifestsChanSize)

	chromedp.ListenTarget(ctx, createNetworkEventHandler(manifests))

	if err := chromedp.Run(ctx, chromedp.Navigate(pageURL)); err != nil {
		return false, err
	}

	select {
	case manifest := <-manifests:
		s.url = manifest
		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (s *HLSStrategy) Download() {}
