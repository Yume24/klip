package strategy

import (
	"context"
	"errors"
	"sync"
	"time"
)

const timeout = 15 * time.Second
const noSuitableStrategyErrorMessage = "cannot download from this site"

type DownloadStrategy interface {
	Scout(ctx context.Context, pageURL string) (bool, error)
	Download()
}

func iterateStrategies(strategies []DownloadStrategy, pageURL string, b *browser, ch chan<- DownloadStrategy) {
	var wg sync.WaitGroup

	for _, strategy := range strategies {
		wg.Go(func() {
			ctx, cleanup := b.createNewBrowserContextWithTimeout(timeout)
			defer cleanup()

			canHandle, err := strategy.Scout(ctx, pageURL)
			if err != nil {
				return
			}

			if canHandle {
				ch <- strategy
			}
		})
	}

	go func() { wg.Wait(); close(ch) }()

}

func GetDownloadStrategy(pageURL string, strategies []DownloadStrategy) (DownloadStrategy, error) {
	browser := &browser{}
	_, cleanup := browser.createNewBrowserContext()
	defer cleanup()

	ch := make(chan DownloadStrategy, len(strategies))

	iterateStrategies(strategies, pageURL, browser, ch)

	if strategy, ok := <-ch; ok {
		return strategy, nil
	}

	return nil, errors.New(noSuitableStrategyErrorMessage)
}
