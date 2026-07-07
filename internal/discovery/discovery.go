package discovery

import (
	"net/url"

	"github.com/chromedp/cdproto/network"
)

const manifestsChanSize = 1

func createEventHandler(manifests chan<- *url.URL) networkEventHandler {
	return func(event *network.EventResponseReceived) {
		if isMediaManifest(event) {
			parsedUrl, err := url.Parse(event.Response.URL)
			if err != nil {
				return
			}

			select {
			case manifests <- parsedUrl:
			default:
			}
		}
	}
}

// Returns the URL pointing to video manifest
func DiscoverManifestURL(pageURL string, discoverer Discoverer) (*url.URL, error) {
	manifests := make(chan *url.URL, manifestsChanSize)

	browserCtx, cleanup, err := initializeBrowser(discoverer.isHeadless(), createEventHandler(manifests))
	if err != nil {
		return nil, err
	}

	defer cleanup()

	return discoverer.discoverMediaManifest(browserCtx, pageURL, manifests)
}
