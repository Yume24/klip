package discovery

import (
	"context"
	"net/url"
)

type Discoverer interface {
	discoverMediaManifest(ctx context.Context, pageURL string, manifests <-chan *url.URL) (*url.URL, error)
	isHeadless() bool
}
