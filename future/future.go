package future

import (
	"context"
	"errors"
	"sync"
)

type Future[T any] interface {
	Task[T]
	Execute()
	Done() <-chan struct{}
	IsDone() bool
}

func New[T any](ctx context.Context, task Task[T]) Future[T] {
	return newFuture(ctx, task)
}

func Execute[T any](ctx context.Context, task Task[T]) Future[T] {
	f := newFuture(ctx, task)
	f.Execute()
	return f
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
	result T
	err    error
}

func (f *future[T]) Execute() {
	go f.once.Do(func() {
		f.result, f.err = f.task.Result(f.ctx)
		f.cancel(errDone)
	})
}

func (f *future[T]) Result(ctx context.Context) (T, error) {
	f.Execute()
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
