package results

import (
	"fmt"

	"github.com/cmgsj/go-lib/options"
)

func Ok[T any, E error](value T) Result[T, E] {
	return &resultOk[T, E]{value: value}
}

type resultOk[T any, E error] struct {
	value T
}

func (r *resultOk[T, E]) Ok() options.Option[T] {
	return options.Some(r.value)
}

func (r *resultOk[T, E]) Err() options.Option[E] {
	return options.None[E]()
}

func (r *resultOk[T, E]) IsOk() bool {
	return true
}

func (r *resultOk[T, E]) IsErr() bool {
	return false
}

func (r *resultOk[T, E]) Expect(message string) T {
	return r.value
}

func (r *resultOk[T, E]) ExpectErr(message string) E {
	panic(message + ": " + fmt.Sprint(r.value))
}

func (r *resultOk[T, E]) Unwrap() T {
	return r.value
}

func (r *resultOk[T, E]) UnwrapErr() E {
	panic("called Result.UnwrapErr on an ResultOk value: " + fmt.Sprint(r.value))
}

func (r *resultOk[T, E]) result() {}
