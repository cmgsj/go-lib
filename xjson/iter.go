package xjson

import (
	"iter"
	"strconv"
)

const Root string = "."

func All(v *Value) iter.Seq2[string, *Value] {
	return func(yield func(string, *Value) bool) {
		iterate(Root, v, yield)
	}
}

func Keys(v *Value) iter.Seq[string] {
	return func(yield func(string) bool) {
		iterate(Root, v, func(key string, _ *Value) bool {
			return yield(key)
		})
	}
}

func Values(v *Value) iter.Seq[*Value] {
	return func(yield func(*Value) bool) {
		iterate(Root, v, func(_ string, value *Value) bool {
			return yield(value)
		})
	}
}

func iterate(key string, value *Value, yield func(key string, value *Value) bool) bool {
	if !yield(key, value) {
		return false
	}

	switch value.kind {
	case KindArray:
		for i, v := range value.GetArray() {
			k := key + `[` + strconv.Itoa(i) + `]`
			if !iterate(k, v, yield) {
				return false
			}
		}

	case KindObject:
		for k, v := range value.GetObject() {
			k := key + `["` + k + `"]`
			if !iterate(k, v, yield) {
				return false
			}
		}
	}

	return true
}
