package future

import "context"

func Value[T any](value T) Future[T] {
	return Result(value, nil)
}

func Error[T any](err error) Future[T] {
	var v T
	return Result(v, err)
}

func Result[T any](result T, err error) Future[T] {
	f := &resultFuture[T]{
		result: result,
		err:    err,
		done:   make(chan struct{}),
	}

	close(f.done)

	return f
}

type resultFuture[T any] struct {
	result T
	err    error
	done   chan struct{}
}

func (f *resultFuture[T]) Get(ctx context.Context) (T, error) {
	var v T

	select {
	case <-ctx.Done():
		return v, ctx.Err()

	default:
		if f.err != nil {
			return v, f.err
		}

		return f.result, nil
	}
}

func (f *resultFuture[T]) Done() <-chan struct{} {
	return f.done
}

func (f *resultFuture[T]) IsDone() bool {
	return true
}
