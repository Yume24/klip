package hls

import (
	"net/url"
	"strings"

	"github.com/chromedp/cdproto/network"
)

const manifestExtension = ".m3u8"

var relevantContentTypes = map[string]bool{
	"application/vnd.apple.mpegurl": true,
	"application/x-mpegurl":         true,
	"application/mpegurl":           true,
}

func hasManifestContentType(contentType string) bool {
	mime := strings.ToLower(strings.TrimSpace(contentType))
	return relevantContentTypes[mime]
}

func hasManifestURLSuffix(rawURL string) bool {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return strings.HasSuffix(strings.ToLower(parsed.Path), manifestExtension)
}

// Predicate for deciding if a given network response is a media manifest
func isMediaManifest(response *network.EventResponseReceived) bool {
	return hasManifestContentType(response.Response.MimeType) || hasManifestURLSuffix(response.Response.URL)
}
