package options

type Option[T any] interface {
	IsSome() bool
	IsNone() bool
	Expect(message string) T
	Unwrap() T
	internal()
}

func New[T any](value T, ok bool) Option[T] {
	if !ok {
		return None[T]()
	}

	return Some(value)
}
