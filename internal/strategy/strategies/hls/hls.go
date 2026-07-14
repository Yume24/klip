package hls

import (
	"bytes"
	"context"

	"github.com/Eyevinn/hls-m3u8/m3u8"
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

func (s *HLSStrategy) Download() error {
	buf := bytes.Buffer{}
	if err := getResponseBody(s.url, &buf); err != nil {
		return err
	}

	playlist, listType, err := m3u8.Decode(buf, true)
	if err != nil {
		return err
	}

	switch listType {
	case m3u8.MEDIA:
		if err := handleMediaPlaylist(playlist.(*m3u8.MediaPlaylist), s.url); err != nil {
			return err
		}
	case m3u8.MASTER:
		if err := handleMasterPlaylist(playlist.(*m3u8.MasterPlaylist), s.url); err != nil {
			return err
		}
	}

	return nil
}
