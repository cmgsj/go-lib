package results

import "github.com/cmgsj/go-lib/options"

type Result[T any, E error] interface {
	Ok() options.Option[T]
	Err() options.Option[E]
	IsOk() bool
	IsErr() bool
	Expect(message string) T
	ExpectErr(message string) E
	Unwrap() T
	UnwrapErr() E
	result()
}

func New[T any](value T, err error) Result[T, error] {
	if err != nil {
		return Err[T](err)
	}

	return Ok[T, error](value)
}
