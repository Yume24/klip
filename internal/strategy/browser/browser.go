package browser

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type Browser struct {
	ctx context.Context
}

func (b *Browser) initializeContextIfEmpty() {
	if b.ctx == nil {
		b.ctx = context.Background()
	}
}

func (b *Browser) CreateNewBrowserContext() (context.Context, context.CancelFunc) {
	b.initializeContextIfEmpty()
	ctx, cleanup := chromedp.NewContext(b.ctx)
	b.ctx = ctx
	return ctx, cleanup
}

func (b *Browser) CreateNewBrowserContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	b.initializeContextIfEmpty()
	timeoutCtx, timeoutCtxCancel := context.WithTimeout(b.ctx, timeout)
	browserCtx, browserCtxCancel := chromedp.NewContext(timeoutCtx)

	cleanup := func() {
		browserCtxCancel()
		timeoutCtxCancel()
	}

	return browserCtx, cleanup
}
