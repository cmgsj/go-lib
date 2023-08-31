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

type Task[T any] func(context.Context) (T, error)

func New[T any](ctx context.Context, task Task[T]) Future[T] {
	ctx, cancel := context.WithCancelCause(ctx)
	f := &future[T]{ctx: ctx}
	go func() {
		f.val, f.err = task(ctx)
		cancel(errDone)
	}()
	return f
}

var errDone = errors.New("future done")

type future[T any] struct {
	ctx context.Context
	val T
	err error
}

func (f *future[T]) Get(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		return f.val, ctx.Err()
	case <-f.ctx.Done():
		if context.Cause(f.ctx) != errDone {
			return f.val, f.ctx.Err()
		}
		return f.val, f.err
	}
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
