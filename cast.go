package xiter

import (
	"iter"
)

// Cast converts a sequence of any into a sequence of E values, returning the
// converted value and a success flag for each element.
func Cast[E any](s iter.Seq[any]) iter.Seq2[E, bool] {
	return func(yield func(E, bool) bool) {
		for a := range s {
			if e, ok := a.(E); !yield(e, ok) {
				return
			}
		}
	}
}
