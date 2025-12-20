package xiter

import (
	"iter"
)

type Pair[K, V any] struct {
	Key   K
	Value V
}

// GroupBy groups sequence elements using the provided key function and returns
// the grouped result as map[K][]E.
func GroupBy[E any, K comparable](s iter.Seq[E], f func(E) K) map[K][]E {
	m := make(map[K][]E)
	for e := range s {
		k := f(e)
		m[k] = append(m[k], e)
	}
	return m
}

// ToSet converts the sequence into a map[E]struct{} set containing all unique
// elements.
func ToSet[E comparable](s iter.Seq[E]) map[E]struct{} {
	set := make(map[E]struct{})
	for e := range s {
		set[e] = struct{}{}
	}
	return set
}

// ToSlice converts the sequence into a []E slice containing all elements.
func ToSlice[E any](s iter.Seq[E]) []E {
	var slice []E
	for e := range s {
		slice = append(slice, e)
	}
	return slice
}

// ToMap converts a key/value sequence into a map[K]V containing all pairs.
func ToMap[K comparable, V any](s iter.Seq2[K, V]) map[K]V {
	m := make(map[K]V)
	for k, v := range s {
		m[k] = v
	}
	return m
}

// ToSlice2 converts a key/value sequence into a []Pair[K, V] slice containing
// all pairs.
func ToSlice2[K, V any](s iter.Seq2[K, V]) []Pair[K, V] {
	var slice []Pair[K, V]
	for k, v := range s {
		slice = append(slice, Pair[K, V]{Key: k, Value: v})
	}
	return slice
}
