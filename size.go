package xiter

import (
	"iter"
)

// Size returns the number of elements in the sequence.
func Size[E any](s iter.Seq[E]) int {
	count := 0
	for range s {
		count++
	}
	return count
}

// Size2 returns the number of key/value pairs in the sequence.
func Size2[K, V any](s iter.Seq2[K, V]) int {
	count := 0
	for range s {
		count++
	}
	return count
}

// SizeFunc returns the number of elements that satisfy the predicate.
func SizeFunc[E any](s iter.Seq[E], f func(E) bool) int {
	count := 0
	for e := range s {
		if f(e) {
			count++
		}
	}
	return count
}

// SizeFunc2 returns the number of key/value pairs that satisfy the predicate.
func SizeFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) int {
	count := 0
	for k, v := range s {
		if f(k, v) {
			count++
		}
	}
	return count
}

// SizeValue returns the number of elements equal to the specified value.
func SizeValue[E comparable](s iter.Seq[E], v E) int {
	return SizeFunc(s, func(e E) bool { return e == v })
}

// SizeValue2 returns the number of key/value pairs whose key and value equal the
// specified values.
func SizeValue2[K, V comparable](s iter.Seq2[K, V], k K, v V) int {
	return SizeFunc2(s, func(ck K, cv V) bool { return ck == k && cv == v })
}
