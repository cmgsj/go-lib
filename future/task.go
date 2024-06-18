package future

import "context"

type Task[T any] interface {
	Execute(ctx context.Context) (T, error)
}

type TaskFunc[T any] func(ctx context.Context) (T, error)

func (f TaskFunc[T]) Execute(ctx context.Context) (T, error) {
	return f(ctx)
}
