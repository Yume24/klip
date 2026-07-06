package discovery

import (
	"net/url"
)

// Carries the relevant information about the network response
type networkResponse struct {
	url         *url.URL
	contentType string
}
