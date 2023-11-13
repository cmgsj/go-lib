package future

import (
	"context"
	"errors"
)

type Future[T any] interface {
	Get(ctx context.Context) (T, error)
	IsReady() bool
	Done() <-chan struct{}
}

type Task[T any] interface {
	Execute(ctx context.Context) (T, error)
}

type TaskFunc[T any] func(ctx context.Context) (T, error)

func (f TaskFunc[T]) Execute(ctx context.Context) (T, error) {
	return f(ctx)
}

func New[T any](ctx context.Context, task Task[T]) Future[T] {
	ctx, cancel := context.WithCancelCause(ctx)
	f := &future[T]{Context: ctx}
	go func() {
		f.val, f.err = task.Execute(ctx)
		cancel(errDone)
	}()
	return f
}

var errDone = errors.New("done")

type future[T any] struct {
	context.Context
	val T
	err error
}

func (f *future[T]) Get(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		return f.val, ctx.Err()
	case <-f.Done():
		if context.Cause(f) != errDone {
			return f.val, f.Err()
		}
		return f.val, f.err
	}
}

func (f *future[T]) IsReady() bool {
	select {
	case <-f.Done():
		return true
	default:
		return false
	}
}
