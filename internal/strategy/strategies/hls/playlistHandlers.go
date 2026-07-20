package hls

import (
	"bytes"
	"errors"
	"fmt"

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
		return handleMediaPlaylist(playlist, playlistURL)
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

	for _, path := range paths {
		fmt.Println(path)
	}

	return nil
}
