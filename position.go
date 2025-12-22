package xiter

import (
	"iter"
)

// Position returns the index of the first element that satisfies the
// predicate.
func Position[E any](s iter.Seq[E], f func(E) bool) (int, bool) {
	index := 0
	for e := range s {
		if f(e) {
			return index, true
		}
		index++
	}
	return -1, false
}

// Position2 returns the index of the first key/value pair that satisfies the
// predicate.
func Position2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) (int, bool) {
	index := 0
	for k, v := range s {
		if f(k, v) {
			return index, true
		}
		index++
	}
	return -1, false
}
