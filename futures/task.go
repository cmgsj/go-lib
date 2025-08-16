package futures

import (
	"context"
	"errors"
)

type TaskFunc[T any] func(ctx context.Context) (T, error)

func Task[T any](ctx context.Context, task TaskFunc[T]) Future[T] {
	cctx, cancel := context.WithCancelCause(ctx)

	f := &taskFuture[T]{
		ctx: cctx,
	}

	go func() {
		f.value, f.err = task(ctx)
		cancel(errTaskDone)
	}()

	return f
}

var errTaskDone = errors.New("task done")

type taskFuture[T any] struct {
	ctx   context.Context
	value T
	err   error
}

func (f *taskFuture[T]) Get(ctx context.Context) (T, error) {
	var value T

	select {
	case <-ctx.Done():
		return value, ctx.Err()

	case <-f.ctx.Done():
		if context.Cause(f.ctx) == errTaskDone {
			if f.err != nil {
				return value, f.err
			}

			return f.value, nil
		}

		return value, f.ctx.Err()
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
