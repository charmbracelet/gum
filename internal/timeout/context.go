// Package timeout handles context timeouts.
package timeout

import (
	"context"
	"time"
)

// Context setup a new context that times out if the given timeout is > 0.
func Context(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if timeout == 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}
