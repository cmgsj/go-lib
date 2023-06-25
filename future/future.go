package future

import (
	"context"
)

type Future[T any] interface {
	Get(context.Context) (T, error)
	IsReady() bool
	Done() <-chan struct{}
}

type Task[T any] func() (T, error)

func New[T any](task Task[T]) Future[T] {
	f := &future[T]{done: make(chan struct{})}
	go func() {
		f.val, f.err = task()
		close(f.done)
	}()
	return f
}

type future[T any] struct {
	val  T
	err  error
	done chan struct{}
}

func (f *future[T]) Get(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		var val T
		return val, ctx.Err()
	case <-f.done:
		return f.val, f.err
	}
}

func (f *future[T]) IsReady() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

func (f *future[T]) Done() <-chan struct{} {
	return f.done
}
