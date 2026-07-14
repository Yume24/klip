package hls

import (
	"io"
	"net/http"
)

func getResponseBody(url string, dest io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(dest, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
