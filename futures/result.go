package futures

import "context"

func Value[T any](value T) Future[T] {
	return Result(value, nil)
}

func Error[T any](err error) Future[T] {
	var value T
	return Result(value, err)
}

func Result[T any](value T, err error) Future[T] {
	f := &resultFuture[T]{
		value: value,
		err:   err,
		done:  make(chan struct{}),
	}

	close(f.done)

	return f
}

type resultFuture[T any] struct {
	value T
	err   error
	done  chan struct{}
}

func (f *resultFuture[T]) Get(ctx context.Context) (T, error) {
	var value T

	select {
	case <-ctx.Done():
		return value, ctx.Err()

	default:
		if f.err != nil {
			return value, f.err
		}

		return f.value, nil
	}
}

func (f *resultFuture[T]) Done() <-chan struct{} {
	return f.done
}

func (f *resultFuture[T]) IsDone() bool {
	return true
}
