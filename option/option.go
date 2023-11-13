package option

type Option[T any] interface {
	Get() T
	Ok() bool
	Or(other T) T
}

func Some[T any](value T) Option[T] {
	return &some[T]{value: value}
}

func None[T any]() Option[T] {
	return &none[T]{}
}

type some[T any] struct{ value T }

func (o *some[T]) Get() T       { return o.value }
func (o *some[T]) Ok() bool     { return true }
func (o *some[T]) Or(other T) T { return o.value }

type none[T any] struct{}

func (o *none[T]) Get() T       { panic("none option") }
func (o *none[T]) Ok() bool     { return false }
func (o *none[T]) Or(other T) T { return other }
