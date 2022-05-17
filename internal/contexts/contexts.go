package contexts

import (
	"context"
	"time"
)

// Timeout is a context with a timeout of a given duration.
func Timeout(duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), duration)
}
