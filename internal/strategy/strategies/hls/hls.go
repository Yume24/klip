package hls

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/Eyevinn/hls-m3u8/m3u8"
	"github.com/Yume24/klip/internal/utils"
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

func (s *HLSStrategy) Download(output string) error {
	buf := bytes.Buffer{}
	if err := utils.GetResponseBody(s.url, &buf); err != nil {
		return err
	}

	playlist, listType, err := m3u8.Decode(buf, true)
	if err != nil {
		return err
	}

	switch listType {
	case m3u8.MEDIA:
		filePath, err := handleMediaPlaylist(playlist.(*m3u8.MediaPlaylist), s.url)
		if err != nil {
			return err
		}
		defer os.Remove(filePath)

		return convertToMP4([]string{filePath}, output)
	case m3u8.MASTER:
		filePaths, err := handleMasterPlaylist(playlist.(*m3u8.MasterPlaylist), s.url)
		if err != nil {
			return err
		}
		for _, filePath := range filePaths {
			defer os.Remove(filePath)
		}

		return convertToMP4(filePaths, output)
	}

	return nil
}

func convertToMP4(pathsToTemp []string, outputPath string) error {
	inputArgs := make([]string, 0)
	for _, path := range pathsToTemp {
		inputArgs = append(inputArgs, "-i", path)
	}

	cmd := exec.Command("ffmpeg", append([]string{"-y"}, append(inputArgs, "-c", "copy", outputPath+".mp4")...)...)
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("ffmpeg", append([]string{"-y"}, append(inputArgs, "-c:v", "libx264", "-c:a", "aac", outputPath+".mp4")...)...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
