package hls

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/Eyevinn/hls-m3u8/m3u8"
)

var errUnsupportedPlaylist = errors.New("unsupported playlist type")

func handleMasterPlaylist(masterPlaylist *m3u8.MasterPlaylist, playlistURL string) error {
	variantBuf := bytes.Buffer{}
	mediaURI, err := resolveAbsoluteURL(playlistURL, masterPlaylist.Variants[0].URI)
	if err != nil {
		return err
	}
	if err := getResponseBody(mediaURI, &variantBuf); err != nil {
		return err
	}

	playlist, _, err := m3u8.Decode(variantBuf, true)
	if err != nil {
		return err
	}
	if playlist, ok := playlist.(*m3u8.MediaPlaylist); ok {
		return handleMediaPlaylist(playlist, mediaURI)
	}

	return errUnsupportedPlaylist
}

func handleMediaPlaylist(playlist *m3u8.MediaPlaylist, playlistURL string) error {
	if !playlist.Closed {
		return errUnsupportedPlaylist
	}

	paths, err := getAllSegments(playlist, playlistURL)
	if err != nil {
		return err
	}

	concatSegments("test.ts", paths)
	deleteSegments(paths)

	return nil
}

func concatSegments(path string, paths []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, segmentFilePath := range paths {
		err := func() error {
			segmentFile, err := os.Open(segmentFilePath)
			if err != nil {
				return err
			}
			defer segmentFile.Close()
			_, err = io.Copy(f, segmentFile)
			if err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteSegments(paths []string) error {
	for _, file := range paths {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}
