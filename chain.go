package xiter

import (
	"iter"
)

// Chain concatenates two sequences into a single sequence.
func Chain[E any](seq1, seq2 iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range seq1 {
			if !yield(e) {
				return
			}
		}
		for e := range seq2 {
			if !yield(e) {
				return
			}
		}
	}
}

// Chain2 concatenates two key/value sequences into a single sequence.
func Chain2[K, V any](seq1, seq2 iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq1 {
			if !yield(k, v) {
				return
			}
		}
		for k, v := range seq2 {
			if !yield(k, v) {
				return
			}
		}
	}
}
