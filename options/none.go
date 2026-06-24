package options

func None[T any]() Option[T] {
	return (*optionNone[T])(nil)
}

type optionNone[T any] struct{}

func (o *optionNone[T]) IsSome() bool {
	return false
}

func (o *optionNone[T]) IsNone() bool {
	return true
}

func (o *optionNone[T]) Expect(message string) T {
	panic(message)
}

func (o *optionNone[T]) Unwrap() T {
	panic("called Option.Unwrap on a OptionNone value")
}

func (o *optionNone[T]) result() {}
