package options

func None[T any]() Option[T] {
	return (*OptionNone[T])(nil)
}

type OptionNone[T any] struct{}

func (o *OptionNone[T]) IsSome() bool {
	return false
}

func (o *OptionNone[T]) IsNone() bool {
	return true
}

func (o *OptionNone[T]) Expect(message string) T {
	panic(message)
}

func (o *OptionNone[T]) Unwrap() T {
	panic("called [Option.Unwrap] on a [OptionNone] value")
}

func (o *OptionNone[T]) internal() {}
