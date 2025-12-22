package xiter

import (
	"iter"
)

// FlatMap applies the provided function to each element, which returns a sequence,
// and then flattens the result into a single sequence.
func FlatMap[E1, E2 any](s iter.Seq[E1], f func(E1) iter.Seq[E2]) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e1 := range s {
			for e2 := range f(e1) {
				if !yield(e2) {
					return
				}
			}
		}
	}
}

// Flatten returns a flattened sequence from a sequence of sequences.
func Flatten[E any](s iter.Seq[iter.Seq[E]]) iter.Seq[E] {
	return FlatMap(s, func(e iter.Seq[E]) iter.Seq[E] { return e })
}
