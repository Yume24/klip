package discovery

import (
	"context"
	"net/url"
)

// Carries the relevant information about the network response
type networkResponse struct {
	url         url.URL
	contentType string
}

// Carries the context used to communicate between the browser and handlers
type browserContext struct {
	ctx        context.Context
	stop       context.CancelFunc
	eventsChan chan networkResponse
}
