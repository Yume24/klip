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
const ffmpegCmd = "ffmpeg"
const mp4Extension = ".mp4"
const yFlag = "-y"
const inputFlag = "-i"

var ffmpegCopyArgs = []string{"-c", "copy"}
var ffmpegRemuxArgs = []string{"-c:v", "libx264", "-c:a", "aac"}

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

		return convertToMP4(output, filePath)
	case m3u8.MASTER:
		filePaths, err := handleMasterPlaylist(playlist.(*m3u8.MasterPlaylist), s.url)
		if err != nil {
			return err
		}
		for _, filePath := range filePaths {
			defer os.Remove(filePath)
		}

		return convertToMP4(output, filePaths...)
	}

	return nil
}

func convertToMP4(outputPath string, pathsToTemp ...string) error {
	baseArgs := createBaseFfmpegArgs(outputPath, pathsToTemp)

	ffmpegCopyArgs := append(baseArgs, ffmpegCopyArgs...)
	ffmpegCopyCmd := exec.Command(ffmpegCmd, ffmpegCopyArgs...)

	if err := ffmpegCopyCmd.Run(); err != nil {
		ffmpedRemuxArgs := append(baseArgs, ffmpegRemuxArgs...)
		ffmpegRemuxCmd := exec.Command(ffmpegCmd, ffmpedRemuxArgs...)
		if err := ffmpegRemuxCmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func createBaseFfmpegArgs(outputPath string, inputFiles []string) []string {
	inputArgs := make([]string, 0, len(inputFiles)*2+2)
	for _, path := range inputFiles {
		inputArgs = append(inputArgs, inputFlag, path)
	}
	inputArgs = append(inputArgs, yFlag, outputPath+mp4Extension)
	return inputArgs
}
