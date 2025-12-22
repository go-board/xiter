package xiter

import (
	"iter"
)

// 测试文件中单独定义的函数
func ToSlice[E any](s iter.Seq[E]) []E {
	var slice []E
	for e := range s {
		slice = append(slice, e)
	}
	return slice
}

func ToMap[K comparable, V any](s iter.Seq2[K, V]) map[K]V {
	m := make(map[K]V)
	for k, v := range s {
		m[k] = v
	}
	return m
}

func FromSlice[E any, S ~[]E](s S) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i, e := range s {
			if !yield(i, e) {
				return
			}
		}
	}
}

func FromMap[K comparable, V any, M ~map[K]V](m M) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}
