package hls

import (
	"bytes"
	"fmt"
	"os"

	"github.com/Eyevinn/hls-m3u8/m3u8"
)

func handleMasterPlaylist(playlist *m3u8.MasterPlaylist, playlistURL string) error {
	variantBuf := bytes.Buffer{}
	mediaURI, err := resolveAbsoluteURL(playlistURL, playlist.Variants[0].URI)
	if err != nil {
		return err
	}
	if err := getResponseBody(mediaURI, &variantBuf); err != nil {
		return err
	}

	mediaPlaylist, err := m3u8.NewMediaPlaylist(0, 500000)
	if err != nil {
		return err
	}

	if err := mediaPlaylist.Decode(variantBuf, true); err != nil {
		return err
	}

	return handleMediaPlaylist(mediaPlaylist, mediaURI)
}

func handleMediaPlaylist(playlist *m3u8.MediaPlaylist, playlistURL string) error {

	if !playlist.Closed {
		return fmt.Errorf("live")
	}
	segmentsList, err := getAllSegments(playlist, playlistURL)
	if err != nil {
		return err
	}

	f, _ := os.Create("test.ts")
	defer f.Close()
	for _, segment := range segmentsList {
		f.Write(segment.data)
	}

	return nil
}
