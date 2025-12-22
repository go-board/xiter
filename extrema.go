package xiter

import (
	"cmp"
	"iter"
)

// Max returns the maximum element in the sequence using cmp.Compare. If the
// sequence is empty, it returns the zero value and false.
func Max[E cmp.Ordered](s iter.Seq[E]) (E, bool) {
	return MaxFunc(s, cmp.Compare[E])
}

// MaxFunc returns the maximum element in the sequence using a custom compare
// function. If the sequence is empty, it returns the zero value and false.
func MaxFunc[E any](s iter.Seq[E], cmp func(E, E) int) (E, bool) {
	it, stop := iter.Pull(s)
	defer stop()
	current, ok := it()
	if !ok {
		return current, false
	}
	for elem, ok := it(); ok; elem, ok = it() {
		if cmp(elem, current) > 0 {
			current = elem
		}
	}
	return current, true
}

// Min returns the minimum element in the sequence using cmp.Compare. If the
// sequence is empty, it returns the zero value and false.
func Min[E cmp.Ordered](s iter.Seq[E]) (E, bool) {
	return MinFunc(s, cmp.Compare[E])
}

// MinFunc returns the minimum element in the sequence using a custom compare
// function. If the sequence is empty, it returns the zero value and false.
func MinFunc[E any](s iter.Seq[E], cmp func(E, E) int) (E, bool) {
	it, stop := iter.Pull(s)
	defer stop()
	current, ok := it()
	if !ok {
		return current, false
	}
	for elem, ok := it(); ok; elem, ok = it() {
		if cmp(elem, current) < 0 {
			current = elem
		}
	}
	return current, true
}
