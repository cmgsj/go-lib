package xiter

import "cmp"

func Compare[V cmp.Ordered](x, y V) int {
	return cmp.Compare(x, y)
}

func CompareKey[K, V any](cmp func(K, K) int) func(K, V, K, V) int {
	return func(xk K, _ V, yk K, _ V) int {
		return cmp(xk, yk)
	}
}

func CompareValue[K, V any](cmp func(V, V) int) func(K, V, K, V) int {
	return func(_ K, xv V, _ K, yv V) int {
		return cmp(xv, yv)
	}
}

func CompareReverse[V any](cmp func(V, V) int) func(V, V) int {
	return func(x, y V) int {
		return cmp(y, x)
	}
}

func CompareReverse2[K, V any](cmp func(K, V, K, V) int) func(K, V, K, V) int {
	return func(xk K, xv V, yk K, yv V) int {
		return cmp(yk, yv, xk, xv)
	}
}
