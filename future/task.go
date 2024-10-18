package future

import (
	"context"
	"errors"
)

func Task[T any](ctx context.Context, task func(ctx context.Context) (T, error)) Future[T] {
	cctx, cancel := context.WithCancelCause(ctx)

	f := &taskFuture[T]{
		ctx: cctx,
	}

	go func() {
		f.result, f.err = task(ctx)
		cancel(errTaskDone)
	}()

	return f
}

var errTaskDone = errors.New("task done")

type taskFuture[T any] struct {
	ctx    context.Context
	result T
	err    error
}

func (f *taskFuture[T]) Get(ctx context.Context) (T, error) {
	var v T

	select {
	case <-ctx.Done():
		return v, ctx.Err()

	case <-f.ctx.Done():
		if context.Cause(f.ctx) == errTaskDone {
			if f.err != nil {
				return v, f.err
			}

			return f.result, nil
		}

		return v, f.ctx.Err()
	}
}

func (f *taskFuture[T]) Done() <-chan struct{} {
	return f.ctx.Done()
}

func (f *taskFuture[T]) IsDone() bool {
	select {
	case <-f.ctx.Done():
		return true

	default:
		return false
	}
}
