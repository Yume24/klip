package hls

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/Eyevinn/hls-m3u8/m3u8"
	"github.com/Yume24/klip/internal/utils"
)

const audioType = "audio"

var errUnsupportedPlaylist = errors.New("unsupported playlist type")
var errNoAudio = errors.New("no audio in playlist")

func handleMasterPlaylist(masterPlaylist *m3u8.MasterPlaylist, playlistURL string) ([]string, error) {
	var wg sync.WaitGroup
	errorCh := make(chan error, 1)
	pathsCh := make(chan string, 2)

	variant := decideVariant(masterPlaylist.Variants)

	audioRelativeURI, err := findAudioPlaylist(variant.Alternatives)
	if err != nil && !errors.Is(err, errNoAudio) {
		return nil, err
	} else {
		audioURI, err := utils.ResolveAbsoluteURL(playlistURL, audioRelativeURI)
		if err != nil {
			return nil, err
		}

		wg.Go(func() {
			playlist, err := getMediaPlaylist(audioURI)
			if err != nil {
				select {
				case errorCh <- err:
				default:
				}

				return
			}

			paths, err := handleMediaPlaylist(playlist, audioURI)
			if err != nil {
				select {
				case errorCh <- err:
				default:
				}

				return
			}

			pathsCh <- paths
		})
	}

	mediaURI, err := utils.ResolveAbsoluteURL(playlistURL, variant.URI)
	if err != nil {
		return nil, err
	}

	wg.Go(func() {
		playlist, err := getMediaPlaylist(mediaURI)
		if err != nil {
			select {
			case errorCh <- err:
			default:
			}

			return
		}

		paths, err := handleMediaPlaylist(playlist, mediaURI)
		if err != nil {
			select {
			case errorCh <- err:
			default:
			}

			return
		}

		pathsCh <- paths
	})

	go func() {
		wg.Wait()
		close(errorCh)
	}()

	if err, ok := <-errorCh; ok {
		return nil, err
	}

	paths := make([]string, 0, 2)
	for path := range pathsCh {
		paths = append(paths, path)
	}

	return paths, nil
}

func getMediaPlaylist(playlistURL string) (*m3u8.MediaPlaylist, error) {
	playlistBuf := bytes.Buffer{}
	if err := utils.GetResponseBody(playlistURL, &playlistBuf); err != nil {
		return nil, err
	}

	playlist, _, err := m3u8.Decode(playlistBuf, true)
	if err != nil {
		return nil, err
	}
	if playlist, ok := playlist.(*m3u8.MediaPlaylist); ok {
		return playlist, nil
	}

	return nil, errUnsupportedPlaylist
}

func handleMediaPlaylist(playlist *m3u8.MediaPlaylist, playlistURL string) (string, error) {
	var filePath string

	if !playlist.Closed {
		return filePath, errUnsupportedPlaylist
	}

	paths, err := getAllSegments(playlist, playlistURL)
	if err != nil {
		return filePath, err
	}

	filePath, err = concatSegments(paths)
	if err != nil {
		return filePath, err
	}

	if err := deleteSegments(paths); err != nil {
		return filePath, err
	}

	return filePath, nil
}

func findAudioPlaylist(alternatives []*m3u8.Alternative) (string, error) {
	for _, alternative := range alternatives {
		if strings.ToLower(alternative.Type) == audioType {
			return alternative.URI, nil
		}
	}

	return "", errNoAudio
}

func concatSegments(paths []string) (string, error) {
	var finalPath string

	f, err := os.CreateTemp("", "")
	if err != nil {
		return finalPath, err
	}
	defer f.Close()

	finalPath = f.Name()

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
			return finalPath, err
		}
	}

	return finalPath, nil
}

func deleteSegments(paths []string) error {
	for _, file := range paths {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func decideVariant(variants []*m3u8.Variant) *m3u8.Variant {
	return variants[0]
}
