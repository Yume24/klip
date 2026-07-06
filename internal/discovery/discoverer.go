package discovery

import "context"

type Discoverer interface {
	discoverMediaManifest(ctx context.Context, pageURL string) error
}
