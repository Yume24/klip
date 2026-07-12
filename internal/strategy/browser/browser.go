package browser

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type Browser struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewBrowser() *Browser {
	ctx, cancel := chromedp.NewContext(context.Background())
	return &Browser{ctx: ctx, cancel: cancel}
}

func (b *Browser) NewTab(timeout time.Duration) (context.Context, context.CancelFunc) {
	timeoutCtx, timeoutCancel := context.WithTimeout(b.ctx, timeout)
	tabCtx, tabCancel := chromedp.NewContext(timeoutCtx)

	cleanup := func() {
		tabCancel()
		timeoutCancel()
	}

	return tabCtx, cleanup
}

func (b *Browser) Close() {
	b.cancel()
}
