package hls

import (
	"strings"

	"github.com/chromedp/cdproto/network"
)

const manifestExtension = ".m3u8"

var relevantContentTypes = map[string]bool{
	"application/vnd.apple.mpegurl": true,
	"application/x-mpegurl":         true,
}

func hasManifestContentType(contentType string) bool {
	return relevantContentTypes[contentType]
}

func hasManifestURLSuffix(url string) bool {
	return strings.HasSuffix(url, manifestExtension)
}

// Predicate for deciding if a given network response is a media manifest
func isMediaManifest(response *network.EventResponseReceived) bool {
	return hasManifestContentType(response.Response.MimeType) || hasManifestURLSuffix(response.Response.URL)
}
