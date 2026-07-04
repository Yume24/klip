package discovery

import (
	"context"
	"fmt"
	"klip/internal/core"
)

func inspectIncomingTraffic(ctx context.Context, ch <-chan networkResponse, output chan<- core.Media) {
	for {
		select {
		case <-ctx.Done():
			return
		case response := <-ch:
			fmt.Println(response.contentType)
			if false {
				output <- core.Media{URL: response.url, Type: response.contentType}
			}
		}
	}
}
