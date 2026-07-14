package hls

import (
	"bytes"
	"io"
	"net/url"
	"slices"
	"sync"

	"github.com/Eyevinn/hls-m3u8/m3u8"
)

const errorChanSize = 1

type orderedMediaSegment struct {
	id   int
	data io.Reader
}

func fetchAllSegments(playlist *m3u8.MediaPlaylist, playlistURL string) (<-chan orderedMediaSegment, error) {
	var wg sync.WaitGroup

	segments := make(chan orderedMediaSegment, playlist.Count())
	errors := make(chan error, errorChanSize)

	for i, segment := range playlist.GetAllSegments() {
		wg.Go(func() {
			orderedSegment, err := obtainOrderedSegment(segment, playlistURL, i)
			if err != nil {
				select {
				case errors <- err:
				default:
				}

				return
			}
			segments <- orderedSegment
		})
	}

	wg.Wait()
	close(segments)
	close(errors)

	if err := <-errors; err != nil {
		return nil, err
	}

	return segments, nil
}

func obtainOrderedSegment(segment *m3u8.MediaSegment, playlistURL string, i int) (orderedMediaSegment, error) {
	segmentURL, err := resolveAbsoluteSegmentURL(playlistURL, segment.URI)
	if err != nil {
		return orderedMediaSegment{}, err
	}

	segmentData := &bytes.Buffer{}
	if err := getResponseBody(segmentURL, segmentData); err != nil {
		return orderedMediaSegment{}, err
	}

	// decryptSegment(segmentData)

	return orderedMediaSegment{id: i, data: segmentData}, nil
}

func resolveAbsoluteSegmentURL(playlistURL string, segmentURL string) (string, error) {
	playlistParsedURL, err := url.ParseRequestURI(playlistURL)
	if err != nil {
		return "", err
	}

	return playlistParsedURL.JoinPath(segmentURL).String(), nil
}

func createSegmentsList(segments <-chan orderedMediaSegment) []orderedMediaSegment {
	segmentsList := make([]orderedMediaSegment, 0)
	for segment := range segments {
		segmentsList = append(segmentsList, segment)
	}
	return segmentsList
}

func sortSegments(segmentsList []orderedMediaSegment) {
	slices.SortFunc(segmentsList, func(segment1 orderedMediaSegment, segment2 orderedMediaSegment) int {
		return segment1.id - segment2.id
	})
}
