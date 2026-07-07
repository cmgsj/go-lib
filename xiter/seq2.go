package xiter

import (
	"iter"
	"slices"
)

type Seq2[K, V any] iter.Seq2[K, V]

func Pull2[K, V any](s Seq2[K, V]) (next func() (K, V, bool), stop func()) {
	return iter.Pull2(iter.Seq2[K, V](s))
}

func Func2[K, V any](f func() (K, V, bool)) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			k, v, ok := f()
			if !ok || !yield(k, v) {
				return
			}
		}
	}
}

func Slice2[S ~[]V, V any](s S) Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range s {
			if !yield(i, v) {
				return
			}
		}
	}
}

func Map2[M ~map[K]V, K comparable, V any](m M) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func Pairs2[K, V any](pairs ...Pair[K, V]) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, p := range pairs {
			if !yield(p.Key, p.Value) {
				return
			}
		}
	}
}

func (s Seq2[K, V]) Filter(f func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

func (s Seq2[K, V]) Map[K2, V2 any](f func(K, V) (K2, V2)) Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s {
			k2, v2 := f(k, v)
			if !yield(k2, v2) {
				return
			}
		}
	}
}

func (s Seq2[K, V]) FlatMap[K2, V2 any](f func(K, V) Seq2[K2, V2]) Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s {
			for k2, v2 := range f(k, v) {
				if !yield(k2, v2) {
					return
				}
			}
		}
	}
}

func (s Seq2[K, V]) Reduce[R any](acc func(R, K, V) R) R {
	var r R
	for k, v := range s {
		r = acc(r, k, v)
	}
	return r
}

func (s Seq2[K, V]) Sorted(cmp func(K, V, K, V) int) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var pairs []Pair[K, V]
		for k, v := range s {
			pairs = append(pairs, Pair[K, V]{Key: k, Value: v})
		}
		slices.SortFunc(pairs, func(x, y Pair[K, V]) int {
			return cmp(x.Key, x.Value, y.Key, y.Value)
		})
		for _, p := range pairs {
			if !yield(p.Key, p.Value) {
				return
			}
		}
	}
}

func (s Seq2[K, V]) Skip(n int) Seq2[K, V] {
	if n < 0 {
		panic("cannot be negative")
	}
	return func(yield func(K, V) bool) {
		i := 0
		for k, v := range s {
			if i < n {
				i++
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

func (s Seq2[K, V]) Take(n int) Seq2[K, V] {
	if n < 0 {
		panic("cannot be negative")
	}
	return func(yield func(K, V) bool) {
		if n == 0 {
			return
		}
		i := 0
		for k, v := range s {
			i++
			if !yield(k, v) {
				return
			}
			if i >= n {
				return
			}
		}
	}
}

func (s Seq2[K, V]) First() (K, V, bool) {
	for k, v := range s {
		return k, v, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

func (s Seq2[K, V]) Last() (K, V, bool) {
	var lastK K
	var lastV V
	var found bool
	for k, v := range s {
		lastK = k
		lastV = v
		found = true
	}
	return lastK, lastV, found
}

func (s Seq2[K, V]) Min(cmp func(K, V, K, V) int) (K, V, bool) {
	var minK K
	var minV V
	var found bool
	for k, v := range s {
		if !found || cmp(minK, minV, k, v) < 0 {
			minK = k
			minV = v
			found = true
		}
	}
	return minK, minV, found
}

func (s Seq2[K, V]) Max(cmp func(K, V, K, V) int) (K, V, bool) {
	var maxK K
	var maxV V
	var found bool
	for k, v := range s {
		if !found || cmp(maxK, maxV, k, v) > 0 {
			maxK = k
			maxV = v
			found = true
		}
	}
	return maxK, maxV, found
}

func (s Seq2[K, V]) None(f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return false
		}
	}
	return true
}

func (s Seq2[K, V]) Any(f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return true
		}
	}
	return false
}

func (s Seq2[K, V]) All(f func(K, V) bool) bool {
	for k, v := range s {
		if !f(k, v) {
			return false
		}
	}
	return true
}

func (s Seq2[K, V]) Len() int {
	var n int
	for range s {
		n++
	}
	return n
}

func (s Seq2[K, V]) Range(f func(K, V) bool) {
	for k, v := range s {
		if !f(k, v) {
			return
		}
	}
}

func (s Seq2[K, V]) ToSlice[T any](f func(K, V) T) []T {
	var a []T
	for k, v := range s {
		a = append(a, f(k, v))
	}
	return a
}

func (s Seq2[K, V]) ToMap[K2 comparable, V2 any](key func(K, V) K2, value func(K, V) V2) map[K2]V2 {
	m := make(map[K2]V2)
	for k, v := range s {
		m[key(k, v)] = value(k, v)
	}
	return m
}

func (s Seq2[K, V]) Keys() Seq[K] {
	return func(yield func(K) bool) {
		for k := range s {
			if !yield(k) {
				return
			}
		}
	}
}

func (s Seq2[K, V]) Values() Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func (s Seq2[K, V]) Pairs() Seq[Pair[K, V]] {
	return func(yield func(Pair[K, V]) bool) {
		for k, v := range s {
			if !yield(Pair[K, V]{Key: k, Value: v}) {
				return
			}
		}
	}
}
