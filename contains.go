package xiter

import (
	"iter"
)

// Contains reports whether the sequence contains the specified comparable
// element.
func Contains[E comparable](s iter.Seq[E], v E) bool {
	return ContainsFunc(s, func(e E) bool { return e == v })
}

// Contains2 reports whether the key/value sequence contains the specified key
// and value.
func Contains2[K, V comparable](s iter.Seq2[K, V], k K, v V) bool {
	return ContainsFunc2(s, func(ck K, cv V) bool { return ck == k && cv == v })
}

// ContainsFunc reports whether the sequence contains an element that satisfies
// the predicate.
func ContainsFunc[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

func ContainsFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return true
		}
	}
	return false
}

// Any reports whether any element in the sequence satisfies the predicate.
func Any[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

// Any2 reports whether any key/value pair in the sequence satisfies the
// predicate.
func Any2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return true
		}
	}
	return false
}

// All reports whether all elements in the sequence satisfy the predicate.
func All[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if !f(e) {
			return false
		}
	}
	return true
}

// All2 reports whether all key/value pairs in the sequence satisfy the
// predicate.
func All2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if !f(k, v) {
			return false
		}
	}
	return true
}
