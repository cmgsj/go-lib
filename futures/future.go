package futures

import "context"

type Future[T any] interface {
	Get(ctx context.Context) (T, error)
	Done() <-chan struct{}
	IsDone() bool
}
