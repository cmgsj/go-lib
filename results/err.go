package results

import (
	"fmt"

	"github.com/cmgsj/go-lib/options"
)

func Err[T any, E error](err E) Result[T, E] {
	return &resultErr[T, E]{err: err}
}

type resultErr[T any, E error] struct {
	err E
}

func (r *resultErr[T, E]) Ok() options.Option[T] {
	return options.None[T]()
}

func (r *resultErr[T, E]) Err() options.Option[E] {
	return options.Some(r.err)
}

func (r *resultErr[T, E]) IsOk() bool {
	return false
}

func (r *resultErr[T, E]) IsErr() bool {
	return true
}

func (r *resultErr[T, E]) Expect(message string) T {
	panic(message + ": " + fmt.Sprint(r.err))
}

func (r *resultErr[T, E]) ExpectErr(message string) E {
	return r.err
}

func (r *resultErr[T, E]) Unwrap() T {
	panic("called Result.Unwrap on an ResultErr value: " + fmt.Sprint(r.err))
}

func (r *resultErr[T, E]) UnwrapErr() E {
	return r.err
}

func (r *resultErr[T, E]) result() {}
