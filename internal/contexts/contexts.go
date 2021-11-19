package contexts

import (
	"context"
	"time"
)

// DBContext is a context with a timeout of 3s used for database operations.
func DBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}
