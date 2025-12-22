package xiter

import (
	"iter"
)

// Enumerate returns a sequence of (index, element) tuples, adding an index for
// each element in the input sequence.
func Enumerate[E any](s iter.Seq[E]) iter.Seq2[int, E] {
	i := -1
	return func(yield func(int, E) bool) {
		for e := range s {
			i++
			if !yield(i, e) {
				return
			}
		}
	}
}
