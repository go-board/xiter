package xiter

import (
	"iter"
)

// Take returns the first n elements of the sequence.
func Take[E any](s iter.Seq[E], n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if n <= 0 || !yield(e) {
				return
			}
			n--
		}
	}
}

// Take2 returns the first n elements of a key/value sequence.
func Take2[K, V any](s iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if n <= 0 || !yield(k, v) {
				return
			}
			n--
		}
	}
}

// TakeWhile returns a sequence that yields elements until the predicate fails.
func TakeWhile[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if !f(e) || !yield(e) {
				return
			}
		}
	}
}

// TakeWhile2 returns a key/value sequence that yields elements until the
// predicate fails.
func TakeWhile2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if !f(k, v) || !yield(k, v) {
				return
			}
		}
	}
}

// Skip returns a sequence that skips the first n elements.
func Skip[E any](s iter.Seq[E], n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if n > 0 {
				n--
				continue
			}
			if !yield(e) {
				return
			}
		}
	}
}

// Skip2 returns a key/value sequence that skips the first n elements.
func Skip2[K, V any](s iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if n > 0 {
				n--
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// SkipWhile returns a sequence that skips elements while the predicate holds.
func SkipWhile[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		skip := true
		for e := range s {
			if skip {
				skip = f(e)
				if skip {
					continue
				}
			}
			if !yield(e) {
				return
			}
		}
	}
}

// SkipWhile2 returns a key/value sequence that skips elements while the
// predicate holds.
func SkipWhile2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		skip := true
		for k, v := range s {
			if skip {
				skip = f(k, v)
				if skip {
					continue
				}
			}
			if !yield(k, v) {
				return
			}
		}
	}
}
