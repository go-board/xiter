package xiter

import (
	"iter"
)

// Range1 generates an integer sequence from 0 to end-1 (end not included).
func Range1[N integral](end N) iter.Seq[N] {
	return func(yield func(N) bool) {
		for i := N(0); i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Range2 generates an integer sequence from start to end-1 (end not included).
func Range2[N integral](start, end N) iter.Seq[N] {
	return func(yield func(N) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Range3 generates an integer sequence from start to end-1 with a step size
// (end not included).
func Range3[N integral](start, end, step N) iter.Seq[N] {
	// 如果步长为0，返回空序列
	if step == 0 {
		return Empty[N]()
	}
	// 如果步长为正，但start >= end，返回空序列
	if step > 0 && start >= end {
		return Empty[N]()
	}
	// 如果步长为负，但start <= end，返回空序列
	if step < 0 && start <= end {
		return Empty[N]()
	}

	return func(yield func(N) bool) {
		if step > 0 {
			for i := start; i < end; i += step {
				if !yield(i) {
					return
				}
			}
		} else {
			for i := start; i > end; i += step {
				if !yield(i) {
					return
				}
			}
		}
	}
}

// FromFunc generates a sequence from a function that returns (element,
// continue). The sequence ends when continue is false.
func FromFunc[E any](f func() (E, bool)) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			e, ok := f()
			if !ok {
				return
			}
			if !yield(e) {
				return
			}
		}
	}
}

// FromFunc2 generates a key/value sequence from a function that returns
// (key, value, continue). The sequence ends when continue is false.
func FromFunc2[K, V any](f func() (K, V, bool)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			k, v, ok := f()
			if !ok {
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// Once generates a sequence containing a single element.
func Once[E any](e E) iter.Seq[E] {
	return func(yield func(E) bool) {
		yield(e)
	}
}

// Once2 generates a sequence containing a single key/value pair.
func Once2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(k, v)
	}
}

// Empty generates an empty sequence.
func Empty[E any]() iter.Seq[E] {
	return func(yield func(E) bool) {
	}
}

// Empty2 generates an empty key/value sequence.
func Empty2[K, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
	}
}

// Repeat generates an infinite sequence repeating a single element.
func Repeat[E any](e E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			if !yield(e) {
				return
			}
		}
	}
}

// Repeat2 generates an infinite sequence repeating a single key/value pair.
func Repeat2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			if !yield(k, v) {
				return
			}
		}
	}
}
