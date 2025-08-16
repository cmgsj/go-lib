package futures

import (
	"context"
	"errors"
)

func All[T any](ctx context.Context, futures ...Future[T]) Future[[]T] {
	if len(futures) == 0 {
		return Value[[]T](nil)
	}

	return Task(ctx, func(ctx context.Context) ([]T, error) {
		var results []T
		var errs []error

		for _, f := range futures {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()

			default:
				result, err := f.Get(ctx)
				if err != nil {
					errs = append(errs, err)
				} else {
					results = append(results, result)
				}
			}
		}

		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}

		return results, nil
	})
}

func First[T any](ctx context.Context, first, second Future[T]) Future[T] {
	return Task(ctx, func(ctx context.Context) (T, error) {
		select {
		case <-first.Done():
			return first.Get(ctx)

		case <-second.Done():
			return second.Get(ctx)
		}
	})
}

func Last[T any](ctx context.Context, first, second Future[T]) Future[T] {
	return Task(ctx, func(ctx context.Context) (T, error) {
		select {
		case <-first.Done():
			return second.Get(ctx)

		case <-second.Done():
			return first.Get(ctx)
		}
	})
}
