package options

func Some[T any](value T) Option[T] {
	return &optionSome[T]{value: value}
}

type optionSome[T any] struct{ value T }

func (o *optionSome[T]) IsSome() bool {
	return true
}

func (o *optionSome[T]) IsNone() bool {
	return false
}

func (o *optionSome[T]) Expect(message string) T {
	return o.value
}

func (o *optionSome[T]) Unwrap() T {
	return o.value
}

func (o *optionSome[T]) result() {}
