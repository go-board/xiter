package xiter

import (
	"iter"
)

// Map applies the provided function to each element and returns a new sequence
// with the results.
func Map[E1, E2 any](s iter.Seq[E1], f func(E1) E2) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e := range s {
			if !yield(f(e)) {
				return
			}
		}
	}
}

// Map2 applies the provided function to each element of a key/value sequence
// and returns a new key/value sequence with the results.
func Map2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// MapWhile applies the provided function to each element and returns a new sequence
// with the results, until the function returns false.
func MapWhile[E1, E2 any](s iter.Seq[E1], f func(E1) (E2, bool)) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e1 := range s {
			if e2, ok := f(e1); !ok || !yield(e2) {
				return
			}
		}
	}
}

// MapWhile2 applies the provided function to each element of a key/value sequence
// and returns a new key/value sequence with the results, until the function returns false.
func MapWhile2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2, bool)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range s {
			if k2, v2, ok := f(k1, v1); !ok || !yield(k2, v2) {
				return
			}
		}
	}
}
