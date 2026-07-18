package hls

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
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
	var currentKey [16]byte
	var currentIV [16]byte

	segments := make(chan orderedMediaSegment, playlist.Count())
	errorsCh := make(chan error, errorChanSize)

	for i, segment := range playlist.GetAllSegments() {
		if len(segment.Keys) > 0 {
			segmentKey, err := getAesEncryptionScheme(segment.Keys)
			if err != nil && !errors.Is(err, errNoEncryption) {
				return nil, err
			}

			keyURI, err := resolveAbsoluteURL(playlistURL, segmentKey.URI)
			if err != nil {
				return nil, err
			}
			var keyBuf bytes.Buffer
			err = getResponseBody(keyURI, &keyBuf)
			if err != nil {
				return nil, err
			}
			currentKey = [16]byte(keyBuf.Bytes())
			if segmentKey.IV != "" {
				iv, err := hex.DecodeString(strings.TrimPrefix(segmentKey.IV, hexPrefix))
				if err != nil {
					return nil, err
				}
				currentIV = [16]byte(iv)
			} else {
				binary.BigEndian.PutUint64(currentIV[8:], segment.SeqId)
			}

		}
		key := currentKey
		iv := currentIV
		wg.Go(func() {
			orderedSegment, err := obtainOrderedSegment(segment, playlistURL, i, iv[:], key[:])
			if err != nil {
				select {
				case errorsCh <- err:
				default:
				}

				return
			}
			segments <- orderedSegment
		})
	}

	wg.Wait()
	close(segments)
	close(errorsCh)

	for err := range errorsCh {
		if err != nil {
			return nil, err
		}
	}

	segmentsList := createSegmentsList(segments)
	sortSegments(segmentsList)

	return segmentsList, nil
}

func obtainOrderedSegment(segment *m3u8.MediaSegment, playlistURL string, i int, iv []byte, keyData []byte) (orderedMediaSegment, error) {
	result := orderedMediaSegment{id: i, key: keyData, iv: iv}
	segmentURL, err := resolveAbsoluteURL(playlistURL, segment.URI)
	if err != nil {
		return result, err
	}

	segmentData := &bytes.Buffer{}
	if err := getResponseBody(segmentURL, segmentData); err != nil {
		return result, err
	}

	result.data = segmentData.Bytes()

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
