package discovery

import (
	"context"
	"klip/internal/core"
	"net/url"
	"strings"
)

const manifestExtension = ".m3u8"

var relevantContentTypes = map[string]bool{
	"application/vnd.apple.mpegurl": true,
	"application/x-mpegurl ":        true,
}

func checkContentType(contentType string) bool {
	return relevantContentTypes[contentType]
}

func checkURL(URL url.URL) bool {
	return strings.HasSuffix(URL.Path, manifestExtension)
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
