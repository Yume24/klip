package strategy

import (
	"context"
	"errors"
	"klip/internal/strategy/browser"
	"sync"
	"time"
)

const timeout = 15 * time.Second
const noSuitableStrategyErrorMessage = "cannot download from this site"

var ErrNoSuitableStrategy = errors.New(noSuitableStrategyErrorMessage)

type DownloadStrategy interface {
	Scout(ctx context.Context, pageURL string) (bool, error)
	Download()
}

func iterateStrategies(strategies []DownloadStrategy, pageURL string, b *browser.Browser, ch chan<- DownloadStrategy) {
	var wg sync.WaitGroup

	for _, strategy := range strategies {
		wg.Go(func() {
			ctx, cleanup := b.NewTab(timeout)
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
	browser := browser.NewBrowser()
	defer browser.Close()

	ch := make(chan DownloadStrategy, len(strategies))

	iterateStrategies(strategies, pageURL, browser, ch)

	if strategy, ok := <-ch; ok {
		return strategy, nil
	}

	return nil, ErrNoSuitableStrategy
}
