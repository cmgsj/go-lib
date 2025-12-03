package results

import (
	"fmt"

	"github.com/cmgsj/go-lib/options"
)

func Err[T any, E error](err E) Result[T, E] {
	return &ResultErr[T, E]{Error: err}
}

type ResultErr[T any, E error] struct {
	Error E
}

func (r *ResultErr[T, E]) Ok() options.Option[T] {
	return options.None[T]()
}

func (r *ResultErr[T, E]) Err() options.Option[E] {
	return options.Some(r.Error)
}

func (r *ResultErr[T, E]) IsOk() bool {
	return false
}

func (r *ResultErr[T, E]) IsErr() bool {
	return true
}

func (r *ResultErr[T, E]) Expect(message string) T {
	panic(message + ": " + fmt.Sprint(r.Error))
}

func (r *ResultErr[T, E]) ExpectErr(message string) E {
	return r.Error
}

func (r *ResultErr[T, E]) Unwrap() T {
	panic("called [Result.Unwrap] on an [ResultErr] value: " + fmt.Sprint(r.Error))
}

func (r *ResultErr[T, E]) UnwrapErr() E {
	return r.Error
}

func (r *ResultErr[T, E]) internal() {}
