package strategy

import (
	"context"
	"errors"
	"klip/internal/strategy/browser"
	"sync"
	"time"
)

const timeout = 15 * time.Second
const smokeTestErrorChanSize = 1
const noSuitableStrategyErrorMessage = "cannot download from this site"

var ErrNoSuitableStrategy = errors.New(noSuitableStrategyErrorMessage)

type DownloadStrategy interface {
	Scout(ctx context.Context, pageURL string) bool
	Download()
}

func iterateStrategies(strategies []DownloadStrategy, pageURL string, b *browser.Browser, suitableStrategies chan<- DownloadStrategy) {
	var wg sync.WaitGroup

	for _, strategy := range strategies {
		wg.Go(func() {
			ctx, cleanup := b.NewTab(timeout)
			defer cleanup()

			canHandle := strategy.Scout(ctx, pageURL)

			if canHandle {
				suitableStrategies <- strategy
			}
		})
	}

	go func() { wg.Wait(); close(suitableStrategies) }()

}

func GetDownloadStrategy(pageURL string, strategies []DownloadStrategy) (DownloadStrategy, error) {
	browser := browser.NewBrowser()
	defer browser.Close()

	errors := make(chan error, smokeTestErrorChanSize)
	go smokeTest(browser, pageURL, errors)

	suitableStrategies := make(chan DownloadStrategy, len(strategies))
	iterateStrategies(strategies, pageURL, browser, suitableStrategies)

	select {
	case err := <-errors:
		return nil, err
	case strategy, ok := <-suitableStrategies:
		if ok {
			return strategy, nil
		}

		return nil, ErrNoSuitableStrategy
	}
}
