package options

func Some[T any](value T) Option[T] {
	return &OptionSome[T]{Value: value}
}

type OptionSome[T any] struct{ Value T }

func (o *OptionSome[T]) IsSome() bool {
	return true
}

func (o *OptionSome[T]) IsNone() bool {
	return false
}

func (o *OptionSome[T]) Expect(message string) T {
	return o.Value
}

func (o *OptionSome[T]) Unwrap() T {
	return o.Value
}

func (o *OptionSome[T]) internal() {}
