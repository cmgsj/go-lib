package result

import (
	"fmt"
)

type Result[T any] interface {
	Get() T
	Ok() bool
	Err() error
}

func New[T any](value T, err error) Result[T] {
	if err != nil {
		return &errResult[T]{err: err}
	}
	return &okResult[T]{value: value}
}

func Ok[T any](value T) Result[T] {
	return &okResult[T]{value: value}
}

func Err[T any](err error) Result[T] {
	return &errResult[T]{err: err}
}

type okResult[T any] struct{ value T }

func (r *okResult[T]) Get() T     { return r.value }
func (r *okResult[T]) Ok() bool   { return true }
func (r *okResult[T]) Err() error { return nil }

type errResult[T any] struct{ err error }

func (r *errResult[T]) Get() T     { panic(fmt.Errorf("error result: %w", r.err)) }
func (r *errResult[T]) Ok() bool   { return false }
func (r *errResult[T]) Err() error { return r.err }
