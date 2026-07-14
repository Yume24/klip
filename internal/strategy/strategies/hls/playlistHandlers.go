package hls

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Eyevinn/hls-m3u8/m3u8"
)

type orderedMediaSegment struct {
	id   int
	data io.Reader
}

func handleMasterPlaylist(playlist *m3u8.MasterPlaylist, playlistURL string) error {
	variantBuf := bytes.Buffer{}
	if err := getResponseBody(playlist.Variants[0].URI, &variantBuf); err != nil {
		return err
	}

	mediaPlaylist, err := m3u8.NewMediaPlaylist(0, 500000)
	if err != nil {
		return err
	}

	if err := mediaPlaylist.Decode(variantBuf, true); err != nil {
		return err
	}

	return handleMediaPlaylist(mediaPlaylist, playlistURL)
}

func handleMediaPlaylist(playlist *m3u8.MediaPlaylist, playlistURL string) error {

	segments, err := fetchAllSegments(playlist, playlistURL)
	if err != nil {
		return err
	}

	segmentsList := createSegmentsList(segments)

	sortSegments(segmentsList)

	for _, segment := range segmentsList {
		fmt.Println(segment.id)
	}

	return nil
}
