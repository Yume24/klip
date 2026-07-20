package hls

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func getResponseBody(url string, dest io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got %d response", resp.StatusCode)
	}

	_, err = io.Copy(dest, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func resolveAbsoluteURL(baseURL string, relativeURL string) (string, error) {
	parsedBase, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}

	parsedRelative, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}

	return parsedBase.ResolveReference(parsedRelative).String(), nil
}
