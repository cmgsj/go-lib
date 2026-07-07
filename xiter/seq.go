package xiter

import (
	"iter"
	"slices"
)

type Seq[V any] iter.Seq[V]

func Pull[V any](s Seq[V]) (next func() (V, bool), stop func()) {
	return iter.Pull(iter.Seq[V](s))
}

func Func[V any](f func() (V, bool)) Seq[V] {
	return func(yield func(V) bool) {
		for {
			v, ok := f()
			if !ok || !yield(v) {
				return
			}
		}
	}
}

func Slice[S ~[]V, V any](s S) Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func MapKeys[M ~map[K]V, K comparable, V any](m M) Seq[K] {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

func MapValues[M ~map[K]V, K comparable, V any](m M) Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m {
			if !yield(v) {
				return
			}
		}
	}
}

func PairKeys[K, V any](pairs ...Pair[K, V]) Seq[K] {
	return func(yield func(K) bool) {
		for _, p := range pairs {
			if !yield(p.Key) {
				return
			}
		}
	}
}

func PairValues[K, V any](pairs ...Pair[K, V]) Seq[V] {
	return func(yield func(V) bool) {
		for _, p := range pairs {
			if !yield(p.Value) {
				return
			}
		}
	}
}

func (s Seq[V]) Filter(f func(V) bool) Seq[V] {
	return func(yield func(V) bool) {
		for v := range s {
			if f(v) && !yield(v) {
				return
			}
		}
	}
}

func (s Seq[V]) Map[V2 any](f func(V) V2) Seq[V2] {
	return func(yield func(V2) bool) {
		for v := range s {
			v2 := f(v)
			if !yield(v2) {
				return
			}
		}
	}
}

func (s Seq[V]) FlatMap[V2 any](f func(V) Seq[V2]) Seq[V2] {
	return func(yield func(V2) bool) {
		for v := range s {
			for v2 := range f(v) {
				if !yield(v2) {
					return
				}
			}
		}
	}
}

func (s Seq[V]) Sorted(compare func(V, V) int) Seq[V] {
	return func(yield func(V) bool) {
		var values []V
		for v := range s {
			values = append(values, v)
		}
		slices.SortFunc(values, func(x, y V) int {
			return compare(x, y)
		})
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

func (s Seq[V]) Skip(n int) Seq[V] {
	if n < 0 {
		panic("cannot be negative")
	}
	return func(yield func(V) bool) {
		i := 0
		for v := range s {
			if i < n {
				i++
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

func (s Seq[V]) Take(n int) Seq[V] {
	if n < 0 {
		panic("cannot be negative")
	}
	return func(yield func(V) bool) {
		if n == 0 {
			return
		}
		i := 0
		for v := range s {
			i++
			if !yield(v) {
				return
			}
			if i >= n {
				return
			}
		}
	}
}

func (s Seq[V]) First() (V, bool) {
	for v := range s {
		return v, true
	}
	var zeroV V
	return zeroV, false
}

func (s Seq[V]) Last() (V, bool) {
	var lastV V
	var found bool
	for v := range s {
		lastV = v
		found = true
	}
	return lastV, found
}

func (s Seq[V]) Min(compare func(V, V) int) (V, bool) {
	var minV V
	var found bool
	for v := range s {
		if !found || compare(minV, v) < 0 {
			minV = v
			found = true
		}
	}
	return minV, found
}

func (s Seq[V]) Max(compare func(V, V) int) (V, bool) {
	var maxV V
	var found bool
	for v := range s {
		if !found || compare(maxV, v) > 0 {
			maxV = v
			found = true
		}
	}
	return maxV, found
}

func (s Seq[V]) None(f func(V) bool) bool {
	for v := range s {
		if f(v) {
			return false
		}
	}
	return true
}

func (s Seq[V]) Any(f func(V) bool) bool {
	for v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

func (s Seq[V]) All(f func(V) bool) bool {
	for v := range s {
		if !f(v) {
			return false
		}
	}
	return true
}

func (s Seq[V]) Len() int {
	var n int
	for range s {
		n++
	}
	return n
}

func (s Seq[V]) Range(f func(V) bool) {
	for v := range s {
		if !f(v) {
			return
		}
	}
}
