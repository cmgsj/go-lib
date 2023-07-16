package future

import (
	"context"
	"errors"
)

type Future[T any] interface {
	Get() (T, error)
	IsReady() bool
	Done() <-chan struct{}
}

type Task[T any] func(context.Context) (T, error)

func New[T any](ctx context.Context, task Task[T]) Future[T] {
	return newFuture(ctx, task)
}

var errDone = errors.New("future done")

func newFuture[T any](ctx context.Context, task Task[T]) *future[T] {
	ctx, cancel := context.WithCancelCause(ctx)
	f := &future[T]{ctx: ctx}
	go func() {
		f.val, f.err = task(ctx)
		cancel(errDone)
	}()
	return f
}

type future[T any] struct {
	ctx context.Context
	val T
	err error
}

func (f *future[T]) Get() (T, error) {
	<-f.ctx.Done()
	err := f.ctx.Err()
	if err != nil && context.Cause(f.ctx) != errDone {
		return f.val, err
	}
	return f.val, f.err
}

func (f *future[T]) IsReady() bool {
	select {
	case <-f.ctx.Done():
		return true
	default:
		return false
	}
}

func (f *future[T]) Done() <-chan struct{} {
	return f.ctx.Done()
}
