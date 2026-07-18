package hls

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/Eyevinn/hls-m3u8/m3u8"
)

const errorChanSize = 1
const hexPrefix = "0x"

type orderedMediaSegment struct {
	id   int
	data []byte
	key  []byte
	iv   []byte
}

func getAllSegments(playlist *m3u8.MediaPlaylist, playlistURL string) ([]orderedMediaSegment, error) {
	var wg sync.WaitGroup
	fmt.Println("Count():", playlist.Count(), "GetAllSegments():", len(playlist.GetAllSegments()))
	keyCache := createKeyCache()

	segments := make(chan orderedMediaSegment, playlist.Count())
	errors := make(chan error, errorChanSize)

	for i, segment := range playlist.GetAllSegments() {
		wg.Go(func() {
			orderedSegment, err := obtainOrderedSegment(segment, playlistURL, i, keyCache)
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

	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	segmentsList := createSegmentsList(segments)
	sortSegments(segmentsList)

	return segmentsList, nil
}

func obtainOrderedSegment(segment *m3u8.MediaSegment, playlistURL string, i int, keys *keyCache) (orderedMediaSegment, error) {
	result := orderedMediaSegment{id: i}
	segmentURL, err := resolveAbsoluteURL(playlistURL, segment.URI)
	if err != nil {
		return result, err
	}

	segmentData := &bytes.Buffer{}
	if err := getResponseBody(segmentURL, segmentData); err != nil {
		return result, err
	}

	result.data = segmentData.Bytes()

	segmentKey, err := getAesEncryptionScheme(segment.Keys)
	if err != nil {
		if errors.Is(err, errNoEncryption) {
			return result, nil
		} else {
			return result, err
		}
	}

	if segmentKey.IV != "" {
		iv, err := hex.DecodeString(strings.TrimPrefix(segmentKey.IV, hexPrefix))
		if err != nil {
			return result, err
		}
		result.iv = iv
	} else {
		result.iv = make([]byte, ivLength)
		binary.BigEndian.PutUint64(result.iv[8:], segment.SeqId)
	}

	keyURI, err := resolveAbsoluteURL(playlistURL, segmentKey.URI)
	if err != nil {
		return result, err
	}

	key, err := keys.getOrFetch(keyURI)
	if err != nil {
		return result, err
	}

	result.key = key

	err = decryptSegment(&result)
	if err != nil {
		return result, err
	}

	return result, nil
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
