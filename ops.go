package xiter

import (
	"cmp"
	"iter"
)

// Take returns the first n elements of the sequence.
func Take[E any](s iter.Seq[E], n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		count := 0
		for e := range s {
			if count >= n {
				return
			}
			if !yield(e) {
				return
			}
			count++
		}
	}
}

// Take2 returns the first n elements of a key/value sequence.
func Take2[K, V any](s iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		count := 0
		for k, v := range s {
			if count >= n {
				return
			}
			if !yield(k, v) {
				return
			}
			count++
		}
	}
}

// Skip returns a sequence that skips the first n elements.
func Skip[E any](s iter.Seq[E], n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		skipped := 0
		for e := range s {
			if skipped < n {
				skipped++
				continue
			}
			if !yield(e) {
				return
			}
		}
	}
}

// Skip2 returns a key/value sequence that skips the first n elements.
func Skip2[K, V any](s iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		skipped := 0
		for k, v := range s {
			if skipped < n {
				skipped++
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

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

// Map applies the provided function to each element and returns a new sequence
// with the results.
func Map[E1, E2 any](s iter.Seq[E1], f func(E1) E2) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e := range s {
			if !yield(f(e)) {
				return
			}
		}
	}
}

// Map2 applies the provided function to each element of a key/value sequence
// and returns a new key/value sequence with the results.
func Map2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Filter returns a new sequence containing only elements that satisfy the
// predicate.
func Filter[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if f(e) && !yield(e) {
				return
			}
		}
	}
}

// Filter2 returns a new key/value sequence containing only pairs that satisfy
// the predicate.
func Filter2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// TakeWhile returns a sequence that yields elements until the predicate fails.
func TakeWhile[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if !f(e) {
				return
			}
			if !yield(e) {
				return
			}
		}
	}
}

// TakeWhile2 returns a key/value sequence that yields elements until the
// predicate fails.
func TakeWhile2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if !f(k, v) {
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// SkipWhile returns a sequence that skips elements while the predicate holds.
func SkipWhile[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		skipping := true
		for e := range s {
			if skipping && f(e) {
				continue
			}
			skipping = false
			if !yield(e) {
				return
			}
		}
	}
}

// SkipWhile2 returns a key/value sequence that skips elements while the
// predicate holds.
func SkipWhile2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		skipping := true
		for k, v := range s {
			if skipping && f(k, v) {
				continue
			}
			skipping = false
			if !yield(k, v) {
				return
			}
		}
	}
}

// Fold reduces a sequence into a single value by applying the given function
// to each element.
func Fold[E any, A any](s iter.Seq[E], init A, f func(A, E) A) A {
	for e := range s {
		init = f(init, e)
	}
	return init
}

// Fold2 reduces a key/value sequence into a single value by applying the given
// function to each pair.
func Fold2[K, V, A any](s iter.Seq2[K, V], init A, f func(A, K, V) A) A {
	for k, v := range s {
		init = f(init, k, v)
	}
	return init
}

// Size returns the number of elements in the sequence.
func Size[E any](s iter.Seq[E]) int {
	return SizeFunc(s, func(_ E) bool { return true })
}

// Size2 returns the number of elements in a key/value sequence.
func Size2[K, V any](s iter.Seq2[K, V]) int {
	return SizeFunc2(s, func(_ K, _ V) bool { return true })
}

// SizeFunc returns the number of elements that satisfy the predicate.
func SizeFunc[E any](s iter.Seq[E], f func(E) bool) int {
	size := 0
	for e := range s {
		if f(e) {
			size++
		}
	}
	return size
}

// SizeFunc2 returns the number of key/value pairs that satisfy the predicate.
func SizeFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) int {
	size := 0
	for k, v := range s {
		if f(k, v) {
			size++
		}
	}
	return size
}

// SizeValue returns the number of elements equal to the specified value.
func SizeValue[E comparable](s iter.Seq[E], v E) int {
	return SizeFunc(s, func(e E) bool { return e == v })
}

// SizeValue2 returns the number of key/value pairs whose value equals the
// specified value.
func SizeValue2[K any, V comparable](s iter.Seq2[K, V], v V) int {
	return SizeFunc2(s, func(_ K, val V) bool { return val == v })
}

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

// ForEach applies the provided function to each element (no return value),
// typically for side effects such as printing or mutating external state.
func ForEach[E any](s iter.Seq[E], f func(E)) {
	for e := range s {
		f(e)
	}
}

// ForEach2 applies the provided function to each key/value pair (no return
// value), typically for side effects.
func ForEach2[K, V any](s iter.Seq2[K, V], f func(K, V)) {
	for k, v := range s {
		f(k, v)
	}
}

// FilterMap applies the provided function to each element, which returns a
// value and a boolean. Only values with a true boolean are kept.
func FilterMap[E1, E2 any](s iter.Seq[E1], f func(E1) (E2, bool)) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e1 := range s {
			if e2, ok := f(e1); ok && !yield(e2) {
				return
			}
		}
	}
}

// FilterMap2 applies the provided function to each key/value pair, which
// returns a key, value, and boolean. Only pairs with a true boolean are kept.
func FilterMap2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2, bool)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range s {
			if k2, v2, ok := f(k1, v1); ok && !yield(k2, v2) {
				return
			}
		}
	}
}

// Keys extracts all keys from a key/value sequence and returns them as a
// sequence.
func Keys[K, V any](s iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range s {
			if !yield(k) {
				return
			}
		}
	}
}

// Values extracts all values from a key/value sequence and returns them as a
// sequence.
func Values[K, V any](s iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

// Concat concatenates multiple sequences into a single sequence.
func Concat[E any](seqs ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, s := range seqs {
			for e := range s {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Concat2 concatenates multiple key/value sequences into a single sequence.
func Concat2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, s := range seqs {
			for k, v := range s {
				if !yield(k, v) {
					return
				}
			}
		}
	}
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

// Equal reports whether two comparable sequences are equal (same length and
// values).
func Equal[E comparable](x, y iter.Seq[E]) bool {
	return EqualFunc(x, y, func(e1, e2 E) bool { return e1 == e2 })
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

// Equal2 reports whether two key/value sequences are equal (comparable keys and
// values must match).
func Equal2[K, V comparable](x, y iter.Seq2[K, V]) bool {
	return EqualFunc2(x, y, func(k1 K, v1 V, k2 K, v2 V) bool { return k1 == k2 && v1 == v2 })
}

// Sum computes the total of all elements in the sequence (element type must
// support addition).
func Sum[E numbric](s iter.Seq[E]) E {
	return Fold(s, E(0), func(a E, e E) E { return a + e })
}

// Join turns a key/value sequence into an element sequence using the provided
// mapping function.
func Join[E, K, V any](s iter.Seq2[K, V], f func(K, V) E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Split turns an element sequence into a key/value sequence using the provided
// split function.
func Split[K, V, E any](s iter.Seq[E], f func(E) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := range s {
			k, v := f(e)
			if !yield(k, v) {
				return
			}
		}
	}
}

// Any reports whether any element in the sequence satisfies the predicate.
func Any[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

// Any2 reports whether any key/value pair in the sequence satisfies the
// predicate.
func Any2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return true
		}
	}
	return false
}

// All reports whether all elements in the sequence satisfy the predicate.
func All[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if !f(e) {
			return false
		}
	}
	return true
}

// All2 reports whether all key/value pairs in the sequence satisfy the
// predicate.
func All2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if !f(k, v) {
			return false
		}
	}
	return true
}

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

// Distinct returns all unique elements in the sequence, preserving the first
// occurrence order.
func Distinct[E comparable](s iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		m := make(map[E]struct{})
		for e := range s {
			if _, ok := m[e]; !ok {
				m[e] = struct{}{}
				if !yield(e) {
					return
				}
			}
		}
	}
}

// ContainsFunc reports whether the sequence contains an element that satisfies
// the predicate.
func ContainsFunc[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

// Contains reports whether the sequence contains the specified comparable
// element.
func Contains[E comparable](s iter.Seq[E], v E) bool {
	return ContainsFunc(s, func(e E) bool { return e == v })
}

func ContainsFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return true
		}
	}
	return false
}

// Contains2 reports whether the key/value sequence contains the specified key
// and value.
func Contains2[K, V comparable](s iter.Seq2[K, V], k K, v V) bool {
	return ContainsFunc2(s, func(ck K, cv V) bool { return ck == k && cv == v })
}

// Find returns the first element that satisfies the predicate.
func Find[E any](s iter.Seq[E], f func(E) bool) (E, bool) {
	for e := range s {
		if f(e) {
			return e, true
		}
	}
	var zero E
	return zero, false
}

// Find2 returns the first key/value pair that satisfies the predicate.
func Find2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) (K, V, bool) {
	for k, v := range s {
		if f(k, v) {
			return k, v, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Position returns the index of the first element that satisfies the
// predicate.
func Position[E any](s iter.Seq[E], f func(E) bool) (int, bool) {
	index := 0
	for e := range s {
		if f(e) {
			return index, true
		}
		index++
	}
	return -1, false
}

// Position2 returns the index of the first key/value pair that satisfies the
// predicate.
func Position2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) (int, bool) {
	index := 0
	for k, v := range s {
		if f(k, v) {
			return index, true
		}
		index++
	}
	return -1, false
}

// IsSortedFunc reports whether the sequence is sorted according to a consistent
// order (ascending or descending) using the provided comparison function.
func IsSortedFunc[E any](s iter.Seq[E], f func(E, E) int) bool {
	it, stop := iter.Pull(s)
	defer stop()
	first, fOk := it()
	second, sOk := it()
	if !fOk || !sOk {
		return true
	}
	prev := second
	initOrder := f(first, second)
	for x, ok := it(); ok; x, ok = it() {
		order := f(prev, x)
		if initOrder == 0 {
			initOrder = order
		}
		if order != initOrder && order != 0 {
			return false
		}
		prev = x
	}
	return true
}

// IsSorted reports whether an ordered sequence is sorted in ascending order.
func IsSorted[E cmp.Ordered](s iter.Seq[E]) bool {
	return IsSortedFunc(s, cmp.Compare[E])
}
