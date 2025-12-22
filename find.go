package xiter

import (
	"iter"
)

// First returns the first element of the sequence. If the sequence is empty, it returns the zero value and false.
func First[E any](s iter.Seq[E]) (E, bool) {
	for e := range s {
		return e, true
	}
	var zero E
	return zero, false
}

// First2 returns the first key/value pair of the sequence. If the sequence is empty, it returns the zero values and false.
func First2[K, V any](s iter.Seq2[K, V]) (K, V, bool) {
	for k, v := range s {
		return k, v, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// FirstFunc returns the first element that satisfies the predicate. If no such element exists, it returns the zero value and false.
func FirstFunc[E any](s iter.Seq[E], f func(E) bool) (E, bool) {
	for e := range s {
		if f(e) {
			return e, true
		}
	}
	var zero E
	return zero, false
}

// FirstFunc2 returns the first key/value pair that satisfies the predicate. If no such pair exists, it returns the zero values and false.
func FirstFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) (K, V, bool) {
	for k, v := range s {
		if f(k, v) {
			return k, v, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Last returns the last element of the sequence. If the sequence is empty, it returns the zero value and false.
func Last[E any](s iter.Seq[E]) (E, bool) {
	var last E
	found := false
	for e := range s {
		last = e
		found = true
	}
	return last, found
}

// Last2 returns the last key/value pair of the sequence. If the sequence is empty, it returns the zero values and false.
func Last2[K, V any](s iter.Seq2[K, V]) (K, V, bool) {
	var lastK K
	var lastV V
	found := false
	for k, v := range s {
		lastK = k
		lastV = v
		found = true
	}
	return lastK, lastV, found
}

// LastFunc returns the last element that satisfies the predicate. If no such element exists, it returns the zero value and false.
func LastFunc[E any](s iter.Seq[E], f func(E) bool) (E, bool) {
	var last E
	found := false
	for e := range s {
		if f(e) {
			last = e
			found = true
		}
	}
	return last, found
}

// LastFunc2 returns the last key/value pair that satisfies the predicate. If no such pair exists, it returns the zero values and false.
func LastFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) (K, V, bool) {
	var lastK K
	var lastV V
	found := false
	for k, v := range s {
		if f(k, v) {
			lastK = k
			lastV = v
			found = true
		}
	}
	return lastK, lastV, found
}
