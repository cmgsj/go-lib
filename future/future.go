package future

import (
	"context"
	"errors"
)

type Future[T any] interface {
	Get(ctx context.Context) (T, error)
	Done() <-chan struct{}
	IsDone() bool
}

func NewFuture[T any](ctx context.Context, task func(ctx context.Context) (T, error)) Future[T] {
	ctx, cancel := context.WithCancelCause(ctx)

	f := &future[T]{
		ctx: ctx,
	}

	go func() {
		f.result, f.err = task(f.ctx)
		cancel(errDone)
	}()

	return f
}

var errDone = errors.New("done")

type future[T any] struct {
	ctx    context.Context
	result T
	err    error
}

func (f *future[T]) Get(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		return f.result, ctx.Err()

	case <-f.ctx.Done():
		if context.Cause(f.ctx) == errDone {
			return f.result, f.err
		}

		return f.result, f.ctx.Err()
	}
}

func (f *future[T]) Done() <-chan struct{} {
	return f.ctx.Done()
}

func (f *future[T]) IsDone() bool {
	select {
	case <-f.Done():
		return true

	default:
		return false
	}
}
