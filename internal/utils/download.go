package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type FetchJob = func() (string, error)

func RunFetchJobs(jobs []FetchJob) ([]string, error) {
	paths := make([]string, len(jobs))
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	for i, job := range jobs {
		wg.Go(func() {
			path, err := job()
			if err != nil {
				select {
				case errCh <- err:
				default:
				}

				return
			}

			paths[i] = path
		})
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err, ok := <-errCh; ok {
		return nil, err
	}

	return paths, nil
}

func GetResponseBody(url string, dest io.Writer) error {
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

func ResolveAbsoluteURL(baseURL string, relativeURL string) (string, error) {
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

func ResolveURLAndDownload(baseURL, relativeURL string) (*bytes.Buffer, error) {
	dataBuf := &bytes.Buffer{}

	dataURL, err := ResolveAbsoluteURL(baseURL, relativeURL)
	if err != nil {
		return nil, err
	}

	if err := GetResponseBody(dataURL, dataBuf); err != nil {
		return nil, err
	}

	return dataBuf, nil
}
