package discovery

import (
	"context"
	"klip/internal/core"
)

func checkContentType(contentType string) bool {
	return false
}

func checkURL(URL string) bool {
	return false
}

// Preidcate for deciding if a given network response is a media manifest
func isMediaManifest(response networkResponse) bool {
	return checkContentType(response.contentType) || checkURL(response.url)
}

func inspectIncomingTraffic(ctx context.Context, ch <-chan networkResponse, output chan<- core.Media) {
	for {
		select {
		case <-ctx.Done():
			return
		case response := <-ch:
			if isMediaManifest(response) {
				output <- core.Media{URL: response.url, Type: response.contentType}
			}
		}
	}
}
