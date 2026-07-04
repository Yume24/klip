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
	"application/x-mpegurl":         true,
}

func hasManifestContentType(contentType string) bool {
	return relevantContentTypes[contentType]
}

func hasManifestURLSuffix(url url.URL) bool {
	return strings.HasSuffix(url.Path, manifestExtension)
}

// Predicate for deciding if a given network response is a media manifest
func isMediaManifest(response networkResponse) bool {
	return hasManifestContentType(response.contentType) || hasManifestURLSuffix(response.url)
}

func inspectIncomingTraffic(ctx context.Context, ch <-chan networkResponse, output chan<- core.Media) {
	for {
		select {
		case <-ctx.Done():
			return
		case response := <-ch:
			if isMediaManifest(response) {
				select {
				case output <- core.Media{URL: &response.url, Type: response.contentType}:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
