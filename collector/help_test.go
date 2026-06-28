package collector

import "iter"

// seqOf returns an iter.Seq[E] yielding the given elements. The sequence
// honors the yield-return-false stop signal. Called with no args it yields an
// empty sequence.
func seqOf[E any](es ...E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, e := range es {
			if !yield(e) {
				return
			}
		}
	}
}

// seq2From returns an iter.Seq2[K, V] pairing keys and values element-wise.
// The two slices must have equal length; empty slices yield an empty sequence.
func seq2From[K, V any](keys []K, values []V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := range keys {
			if !yield(keys[i], values[i]) {
				return
			}
		}
	}
}
