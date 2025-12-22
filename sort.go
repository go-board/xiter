package xiter

import (
	"cmp"
	"iter"
)

// IsSorted checks if the elements of the sequence are sorted in ascending order.
func IsSorted[E cmp.Ordered](s iter.Seq[E]) bool {
	return IsSortedFunc(s, cmp.Compare[E])
}

// IsSortedFunc reports whether the sequence is sorted according to a consistent
// order (ascending or descending) using the provided comparison function.
func IsSortedFunc[E any](s iter.Seq[E], f func(E, E) int) bool {
	it, stop := iter.Pull(s)
	defer stop()
	prev, ok := it()
	if !ok {
		return true
	}
	var initOrder int
	for {
		curr, ok := it()
		if !ok {
			break
		}
		order := f(prev, curr)
		if initOrder == 0 {
			initOrder = order
		}
		if order != initOrder && order != 0 {
			return false
		}
		prev = curr
	}
	return true
}
