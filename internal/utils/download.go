package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const maxConcurrency = 16
const clientTimeout = 30 * time.Second
const maxRetries = 3
const backoff = time.Second

var httpClient = http.Client{Timeout: clientTimeout}

type FetchJob[T any] = func() (T, error)

func RunFetchJobs[T any](jobs []FetchJob[T]) ([]T, error) {
	results := make([]T, len(jobs))
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(maxConcurrency)
	ctx := context.Background()

	for i, job := range jobs {
		wg.Go(func() {
			if err := sem.Acquire(ctx, 1); err != nil {
				return
			}
			defer sem.Release(1)
			result, err := job()
			if err != nil {
				select {
				case <-ctx.Done():
				case errCh <- err:
				default:
				}

				return
			}

			results[i] = result
		})
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err, ok := <-errCh; ok {
		return nil, err
	}

	return results, nil
}

func GetResponseBody(url string, dest io.Writer) (err error) {
	for attempt := range maxRetries {
		err = func() error {
			resp, err := httpClient.Get(url)
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
		}()

		if err == nil {
			return
		}

		time.Sleep(time.Duration(attempt+1) * backoff)
	}

	return
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
