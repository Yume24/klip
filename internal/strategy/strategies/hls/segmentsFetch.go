package hls

import (
	"bytes"
	"os"
	"strconv"
	"sync"

	"github.com/Eyevinn/hls-m3u8/m3u8"
)

const errorChanSize = 1
const connector = "_"
const segmentFilePerm = 0644

func getAllSegments(playlist *m3u8.MediaPlaylist, playlistURL string) ([]string, error) {
	keys, err := getAllKeys(playlist, playlistURL)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	paths := make([]string, playlist.Count())
	errorsCh := make(chan error, errorChanSize)

	for i, segment := range playlist.GetAllSegments() {
		wg.Go(func() {
			segmentPath, err := downloadSegment(segment, playlistURL, i, keys[i])
			if err != nil {
				select {
				case errorsCh <- err:
				default:
				}

				return
			}

			paths[i] = segmentPath
		})
	}

	go func() {
		wg.Wait()
		close(errorsCh)
	}()

	for err := range errorsCh {
		return nil, err
	}

	return paths, nil
}

func downloadSegment(segment *m3u8.MediaSegment, playlistURL string, i int, decryption decrpytionInfo) (string, error) {
	var path string

	segmentBuf := bytes.Buffer{}
	segmentURL, err := resolveAbsoluteURL(playlistURL, segment.URI)
	if err != nil {
		return path, err
	}

	if err := getResponseBody(segmentURL, &segmentBuf); err != nil {
		return path, err
	}

	decryptedSegment, err := decryptSegment(segmentBuf.Bytes(), decryption.key, decryption.iv)
	if err != nil {
		return path, err
	}

	path = createFileName("", i)
	if err := os.WriteFile(path, decryptedSegment, segmentFilePerm); err != nil {
		return path, err
	}

	return path, nil
}

func createFileName(prefix string, id int) string {
	return prefix + connector + strconv.Itoa(id)
}
