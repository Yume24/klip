package discovery

import (
	"context"
	"fmt"
)

func inspectIncomingTraffic(ctx context.Context, ch <-chan networkResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		case response := <-ch:
			fmt.Println(response.contentType)
		}
	}
}
