package xiter

import (
	"iter"
)

// ForEach applies the provided function to each element (no return value),
// typically for side effects such as printing or mutating external state.
func ForEach[E any](s iter.Seq[E], f func(E)) {
	for e := range s {
		f(e)
	}
}

// ForEach2 applies the provided function to each key/value pair (no return
// value), typically for side effects.
func ForEach2[K, V any](s iter.Seq2[K, V], f func(K, V)) {
	for k, v := range s {
		f(k, v)
	}
}

// Fold reduces a sequence into a single value by applying the given function
// to each element.
func Fold[E any, A any](s iter.Seq[E], init A, f func(A, E) A) A {
	for e := range s {
		init = f(init, e)
	}
	return init
}

// Fold2 reduces a key/value sequence into a single value by applying the given
// function to each pair.
func Fold2[K, V, A any](s iter.Seq2[K, V], init A, f func(A, K, V) A) A {
	for k, v := range s {
		init = f(init, k, v)
	}
	return init
}
