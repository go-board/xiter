package xiter

import (
	"iter"
)

// Filter returns a new sequence containing only elements that satisfy the
// predicate.
func Filter[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if f(e) && !yield(e) {
				return
			}
		}
	}
}

// Filter2 returns a new key/value sequence containing only pairs that satisfy
// the predicate.
func Filter2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// FilterMap applies the provided function to each element, which returns a
// value and a boolean. Only values with a true boolean are kept.
func FilterMap[E1, E2 any](s iter.Seq[E1], f func(E1) (E2, bool)) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e1 := range s {
			if e2, ok := f(e1); ok && !yield(e2) {
				return
			}
		}
	}
}

// FilterMap2 applies the provided function to each key/value pair, which
// returns a key, value, and boolean. Only pairs with a true boolean are kept.
func FilterMap2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2, bool)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range s {
			if k2, v2, ok := f(k1, v1); ok && !yield(k2, v2) {
				return
			}
		}
	}
}
