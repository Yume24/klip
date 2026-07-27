package hls

import (
	"bytes"
	"os"
	"sync"

	"github.com/Eyevinn/hls-m3u8/m3u8"
)

const filePattern = "*.part"

func getAllSegments(playlist *m3u8.MediaPlaylist, playlistURL string) ([]string, error) {
	var wg sync.WaitGroup
	hasMap := playlist.Map != nil
	count := playlist.Count()
	if hasMap {
		count += 1
	}
	paths := make([]string, count)
	segmentPaths := paths
	errorsCh := make(chan error, 1)

	if hasMap {
		segmentPaths = paths[1:]
		wg.Go(func() {
			mapPath, err := downloadMap(playlist.Map.URI, playlistURL)
			if err != nil {
				select {
				case errorsCh <- err:
				default:
				}

				return
			}

			paths[0] = mapPath
		})
	}

	keys, err := getAllKeys(playlist, playlistURL)
	if err != nil {
		return nil, err
	}

	for i, segment := range playlist.GetAllSegments() {
		wg.Go(func() {
			segmentPath, err := downloadSegment(segment, playlistURL, keys[i])
			if err != nil {
				select {
				case errorsCh <- err:
				default:
				}

				return
			}

			segmentPaths[i] = segmentPath
		})
	}

	go func() {
		wg.Wait()
		close(errorsCh)
	}()

	if err, ok := <-errorsCh; ok {
		return nil, err
	}

	return paths, nil
}

func downloadSegment(segment *m3u8.MediaSegment, playlistURL string, decryption decrpytionInfo) (string, error) {
	var path string

	segmentBuf, err := resolveURLAndDownload(playlistURL, segment.URI)
	if err != nil {
		return path, err
	}

	decryptedSegment, err := decryptSegment(segmentBuf.Bytes(), decryption.key, decryption.iv)
	if err != nil {
		return path, err
	}

	path, err = createTempFile(decryptedSegment)
	if err != nil {
		return path, err
	}

	return path, nil
}

func downloadMap(mapURL, playlistURL string) (string, error) {
	var path string

	mapData, err := resolveURLAndDownload(playlistURL, mapURL)
	if err != nil {
		return path, err
	}

	path, err = createTempFile(mapData.Bytes())
	if err != nil {
		return path, err
	}

	return path, nil
}

func resolveURLAndDownload(baseURL, relativeURL string) (*bytes.Buffer, error) {
	dataBuf := &bytes.Buffer{}

	dataURL, err := resolveAbsoluteURL(baseURL, relativeURL)
	if err != nil {
		return nil, err
	}

	if err := getResponseBody(dataURL, dataBuf); err != nil {
		return nil, err
	}

	return dataBuf, nil
}

func createTempFile(data []byte) (string, error) {
	var path string

	f, err := os.CreateTemp("", filePattern)
	if err != nil {
		return path, err
	}

	defer f.Close()

	path = f.Name()

	if _, err := f.Write(data); err != nil {
		return path, err
	}

	return path, nil
}
