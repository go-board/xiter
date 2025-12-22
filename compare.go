package xiter

import (
	"cmp"
	"iter"
)

// Compare reports the comparison result between two comparable sequences.
// It returns:
//   - 0 if the sequences are equal
//   - -1 if the first sequence is "less" than the second
//   - 1 if the first sequence is "greater" than the second
func Compare[E cmp.Ordered](x, y iter.Seq[E]) int {
	return CompareFunc(x, y, cmp.Compare[E])
}

// Compare2 reports the comparison result between two comparable key/value sequences.
// It returns:
//   - 0 if the sequences are equal
//   - -1 if the first sequence is "less" than the second
//   - 1 if the first sequence is "greater" than the second
func Compare2[K, V cmp.Ordered](x, y iter.Seq2[K, V]) int {
	return CompareFunc2(x, y, func(k1 K, v1 V, k2 K, v2 V) int {
		if c := cmp.Compare(k1, k2); c != 0 {
			return c
		}
		return cmp.Compare(v1, v2)
	})
}

// CompareFunc compares two sequences using the provided comparison function.
// It returns:
//   - 0 if the sequences are equal
//   - -1 if the first sequence is "less" than the second
//   - 1 if the first sequence is "greater" than the second
func CompareFunc[E any](x, y iter.Seq[E], f func(E, E) int) int {
	it1, stop1 := iter.Pull(x)
	defer stop1()
	it2, stop2 := iter.Pull(y)
	defer stop2()
	for {
		e1, ok1 := it1()
		e2, ok2 := it2()
		if !ok1 && !ok2 {
			return 0
		}
		if !ok1 {
			return -1
		}
		if !ok2 {
			return 1
		}
		if c := f(e1, e2); c != 0 {
			return c
		}
	}
}

// CompareFunc2 compares two key/value sequences using the provided comparison function.
// It returns:
//   - 0 if the sequences are equal
//   - -1 if the first sequence is "less" than the second
//   - 1 if the first sequence is "greater" than the second
func CompareFunc2[K, V any](x, y iter.Seq2[K, V], f func(K, V, K, V) int) int {
	it1, stop1 := iter.Pull2(x)
	defer stop1()
	it2, stop2 := iter.Pull2(y)
	defer stop2()
	for {
		k1, v1, ok1 := it1()
		k2, v2, ok2 := it2()
		if !ok1 && !ok2 {
			return 0
		}
		if !ok1 {
			return -1
		}
		if !ok2 {
			return 1
		}
		if c := f(k1, v1, k2, v2); c != 0 {
			return c
		}
	}
}

// Equal reports whether two comparable sequences are equal (same length and
// values).
func Equal[E comparable](x, y iter.Seq[E]) bool {
	return EqualFunc(x, y, func(e1, e2 E) bool { return e1 == e2 })
}

// Equal2 reports whether two key/value sequences are equal (comparable keys and
// values must match).
func Equal2[K, V comparable](x, y iter.Seq2[K, V]) bool {
	return EqualFunc2(x, y, func(k1 K, v1 V, k2 K, v2 V) bool { return k1 == k2 && v1 == v2 })
}

// EqualFunc compares two sequences using the provided comparison function.
func EqualFunc[E any](x, y iter.Seq[E], f func(E, E) bool) bool {
	it1, stop1 := iter.Pull(x)
	defer stop1()
	it2, stop2 := iter.Pull(y)
	defer stop2()
	for {
		e1, ok1 := it1()
		e2, ok2 := it2()
		if !ok1 && !ok2 {
			return true
		}
		if ok1 != ok2 || !f(e1, e2) {
			return false
		}
	}
}

// EqualFunc2 compares two key/value sequences using the provided comparison
// function.
func EqualFunc2[K, V any](x, y iter.Seq2[K, V], f func(K, V, K, V) bool) bool {
	it1, stop1 := iter.Pull2(x)
	defer stop1()
	it2, stop2 := iter.Pull2(y)
	defer stop2()
	for {
		k1, v1, ok1 := it1()
		k2, v2, ok2 := it2()
		if !ok1 && !ok2 {
			return true
		}
		if ok1 != ok2 || !f(k1, v1, k2, v2) {
			return false
		}
	}
}
