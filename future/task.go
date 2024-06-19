package future

import "context"

type Task[T any] interface {
	Result(ctx context.Context) (T, error)
}

type TaskFunc[T any] func(ctx context.Context) (T, error)

func (f TaskFunc[T]) Result(ctx context.Context) (T, error) {
	return f(ctx)
}
