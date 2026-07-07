package xiter

func Get[V any](v V) V {
	return v
}

func Get2[K, V any](k K, v V) (K, V) {
	return k, v
}

func GetKey[K, V any](k K, _ V) K {
	return k
}

func GetValue[K, V any](_ K, v V) V {
	return v
}

func GetEmpty[V any](_ V) struct{} {
	return struct{}{}
}
