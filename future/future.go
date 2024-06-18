package future

import (
	"context"
	"errors"
	"sync"
)

type Future[T any] interface {
	Task[T]
	Get(ctx context.Context) (T, error)
	Done() <-chan struct{}
	IsDone() bool
}

func Eager[T any](ctx context.Context, task Task[T]) Future[T] {
	f := newFuture(ctx, task)
	f.execute(ctx)
	return f
}

func Lazy[T any](ctx context.Context, task Task[T]) Future[T] {
	return newFuture(ctx, task)
}

func newFuture[T any](ctx context.Context, task Task[T]) *future[T] {
	ctx, cancel := context.WithCancelCause(ctx)
	return &future[T]{
		ctx:    ctx,
		cancel: cancel,
		task:   task,
	}
}

var errDone = errors.New("done")

type future[T any] struct {
	once   sync.Once
	ctx    context.Context
	cancel context.CancelCauseFunc
	task   Task[T]
	val    T
	err    error
}

func (f *future[T]) Execute(ctx context.Context) (T, error) {
	return f.Get(ctx)
}

func (f *future[T]) Get(ctx context.Context) (T, error) {
	f.execute(ctx)
	select {
	case <-ctx.Done():
		return f.val, ctx.Err()
	case <-f.ctx.Done():
		if context.Cause(f.ctx) == errDone {
			return f.val, f.err
		}
		return f.val, f.ctx.Err()
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

func (f *future[T]) execute(ctx context.Context) {
	go f.once.Do(func() {
		f.val, f.err = f.task.Execute(ctx)
		f.cancel(errDone)
	})
}
