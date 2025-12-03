package results

import (
	"fmt"

	"github.com/cmgsj/go-lib/options"
)

func Ok[T any, E error](value T) Result[T, E] {
	return &ResultOk[T, E]{Value: value}
}

type ResultOk[T any, E error] struct {
	Value T
}

func (r *ResultOk[T, E]) Ok() options.Option[T] {
	return options.Some(r.Value)
}

func (r *ResultOk[T, E]) Err() options.Option[E] {
	return options.None[E]()
}

func (r *ResultOk[T, E]) IsOk() bool {
	return true
}

func (r *ResultOk[T, E]) IsErr() bool {
	return false
}

func (r *ResultOk[T, E]) Expect(message string) T {
	return r.Value
}

func (r *ResultOk[T, E]) ExpectErr(message string) E {
	panic(message + ": " + fmt.Sprint(r.Value))
}

func (r *ResultOk[T, E]) Unwrap() T {
	return r.Value
}

func (r *ResultOk[T, E]) UnwrapErr() E {
	panic("called [Result.UnwrapErr] on an [ResultOk] value: " + fmt.Sprint(r.Value))
}

func (r *ResultOk[T, E]) internal() {}
