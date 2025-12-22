package xiter

import (
	"iter"
)

// Join turns a key/value sequence into an element sequence using the provided
// mapping function.
func Join[E, K, V any](s iter.Seq2[K, V], f func(K, V) E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Split turns an element sequence into a key/value sequence using the provided
// split function.
func Split[K, V, E any](s iter.Seq[E], f func(E) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := range s {
			k, v := f(e)
			if !yield(k, v) {
				return
			}
		}
	}
}
