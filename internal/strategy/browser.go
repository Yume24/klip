package strategy

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type browser struct {
	ctx context.Context
}

func (b *browser) initializeContextIfEmpty() {
	if b.ctx == nil {
		b.ctx = context.Background()
	}
}

func (b *browser) createNewBrowserContext() (context.Context, context.CancelFunc) {
	b.initializeContextIfEmpty()
	ctx, cleanup := chromedp.NewContext(b.ctx)
	b.ctx = ctx
	return ctx, cleanup
}

func (b *browser) createNewBrowserContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	b.initializeContextIfEmpty()
	timeoutCtx, timeoutCtxCancel := context.WithTimeout(b.ctx, timeout)
	browserCtx, browserCtxCancel := chromedp.NewContext(timeoutCtx)

	cleanup := func() {
		browserCtxCancel()
		timeoutCtxCancel()
	}

	return browserCtx, cleanup
}
