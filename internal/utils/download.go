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
type RequestOption = func(*http.Request)
type ByteRange struct{ Length, Offset int64 }

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

func WithByteRange(b ByteRange) RequestOption {
	return func(req *http.Request) {
		if b.Length > 0 {
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", b.Offset, b.Offset+b.Length-1))
		}
	}
}

func GetResponseBody(url string, dst io.Writer, opts ...RequestOption) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	for _, opt := range opts {
		opt(req)
	}

	wantStatusCode := http.StatusOK
	if req.Header.Get("Range") != "" {
		wantStatusCode = http.StatusPartialContent
	}

	data, err := retry(func() (data *bytes.Buffer, err error) {
		code, err := fetch(req, data)
		if err != nil {
			return
		}
		if code != wantStatusCode {
			return nil, fmt.Errorf("got: %d status code, want: %d", code, wantStatusCode)
		}

		return
	}, maxRetries, backoff)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, data)
	return err
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

func fetch(req *http.Request, dst io.Writer) (statusCode int, err error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close() //nolint:errcheck

	_, err = io.Copy(dst, resp.Body)
	return
}

func retry[T any](retriable func() (T, error), maxRetries int, backoff time.Duration) (result T, err error) {
	for attempt := range maxRetries {
		result, err = retriable()
		if err == nil {
			return
		}
		if attempt < maxRetries-1 {
			time.Sleep(backoff * time.Duration(attempt+1))
		}
	}

	return
}
