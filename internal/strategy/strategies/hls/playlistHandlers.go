package hls

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/Eyevinn/hls-m3u8/m3u8"
	"github.com/Yume24/klip/internal/utils"
)

var errUnsupportedPlaylist = errors.New("unsupported playlist type")

func handleMasterPlaylist(masterPlaylist *m3u8.MasterPlaylist, playlistURL string, output string) error {
	variantBuf := bytes.Buffer{}
	mediaURI, err := utils.ResolveAbsoluteURL(playlistURL, masterPlaylist.Variants[0].URI)
	if err != nil {
		return err
	}
	if err := utils.GetResponseBody(mediaURI, &variantBuf); err != nil {
		return err
	}

	playlist, _, err := m3u8.Decode(variantBuf, true)
	if err != nil {
		return err
	}
	if playlist, ok := playlist.(*m3u8.MediaPlaylist); ok {
		return handleMediaPlaylist(playlist, mediaURI, output)
	}

	return errUnsupportedPlaylist
}

func handleMediaPlaylist(playlist *m3u8.MediaPlaylist, playlistURL string, output string) error {
	if !playlist.Closed {
		return errUnsupportedPlaylist
	}

	paths, err := getAllSegments(playlist, playlistURL)
	if err != nil {
		return err
	}

	finalPath, err := concatSegments(paths)
	if err != nil {
		return err
	}

	defer os.Remove(finalPath)

	if err := deleteSegments(paths); err != nil {
		return err
	}

	if err := convertToMP4(finalPath, output); err != nil {
		return err
	}

	return nil
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

func convertToMP4(pathToTemp string, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", pathToTemp, "-c", "copy", outputPath+".mp4")
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("ffmpeg", "-y", "-i", pathToTemp, "-c:v", "libx264", "-c:a", "aac", outputPath+".mp4")
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
