package contextdemo

import (
	"context"
	"time"
)

type requestIDKey string

const requestIDContextKey requestIDKey = "request-id"

func WorkWithTimeout(ctx context.Context, work time.Duration) error {
	select {
	case <-time.After(work):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func StreamNumbers(ctx context.Context, limit int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for i := 0; i < limit; i++ {
			select {
			case <-ctx.Done():
				return
			case output <- i:
			}
		}
	}()

	return output
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

func RequestID(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(requestIDContextKey).(string)
	return value, ok
}
