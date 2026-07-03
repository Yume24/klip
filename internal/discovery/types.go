package discovery

import "context"

type networkResponse struct {
	url         string
	contentType string
}

type browserContext struct {
	ctx        context.Context
	cancel     context.CancelFunc
	eventsChan chan networkResponse
}
