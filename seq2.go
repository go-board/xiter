package xiter

import (
	"cmp"
	"iter"
)

// ============================================================================
// Source
// ============================================================================

// FromFunc2 generates a key/value sequence from a supplier function that
// returns (key, value, continue). The sequence ends when continue is false.
// The supplier is called lazily, only as the returned sequence is consumed.
//
//	i := 0
//	FromFunc2(func() (int, int, bool) {
//	    i++
//	    if i > 3 { return 0, 0, false }
//	    return i, i * 10, true
//	})  // yields (1,10), (2,20), (3,30)
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

// Iterate2 generates a key/value sequence starting from (seedK, seedV), where
// each subsequent pair is produced by applying next to the previous one. The
// sequence ends when next returns ok=false; the pair returned alongside
// ok=false is discarded. It is the Seq2 analog of Iterate.
//
//	Iterate2(0, 1, func(k, v int) (int, int, bool) {
//	    if k >= 3 { return 0, 0, false }
//	    return k + 1, v * 2, true
//	})  // yields (0,1), (1,2), (2,4), (3,8)
func Iterate2[K, V any](seedK K, seedV V, next func(K, V) (K, V, bool)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		k, v := seedK, seedV
		if !yield(k, v) {
			return
		}
		for {
			nk, nv, ok := next(k, v)
			if !ok {
				return
			}
			if !yield(nk, nv) {
				return
			}
			k, v = nk, nv
		}
	}
}

// Once2 generates a sequence containing a single key/value pair. The returned
// sequence yields exactly one pair and then ends; yield's return value is
// ignored since there is nothing further to release.
func Once2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(k, v)
	}
}

// Empty2 generates an empty key/value sequence. Iterating it yields nothing
// and returns immediately.
func Empty2[K, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
	}
}

// Repeat2 generates an infinite sequence repeating a single key/value pair.
// Always pair it with a limiting operator such as Take2 to avoid consuming
// forever.
//
//	Take2(Repeat2("k", 1), 3)  // yields ("k",1), ("k",1), ("k",1)
func Repeat2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			if !yield(k, v) {
				return
			}
		}
	}
}

// ============================================================================
// Transform
// ============================================================================

// Map2 applies f to each key/value pair and yields the resulting (K2, V2)
// pairs. The output length matches the input length. Iteration stops as soon
// as yield returns false.
//
//	Map2(Enumerate(seqOf("a", "b")), func(k int, v string) (string, int) {
//	    return v, k
//	})  // yields ("a",0), ("b",1)
func Map2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// MapWhile2 applies f to each key/value pair and yields the results until f
// returns ok=false, at which point the sequence stops. The pair that caused
// the stop is not yielded.
//
//	MapWhile2(Enumerate(Range1(10)), func(k, v int) (int, int, bool) {
//	    if v < 3 { return k, v * 10, true }
//	    return 0, 0, false
//	})  // yields (0,0), (1,10), (2,20)
func MapWhile2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2, bool)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range s {
			if k2, v2, ok := f(k1, v1); !ok || !yield(k2, v2) {
				return
			}
		}
	}
}

// Inspect2 calls f for each key/value pair as it passes through and yields the
// pair unchanged. Useful for debugging or logging in the middle of a lazy
// pipeline. f is called only when the returned sequence is consumed, and only
// for pairs actually visited.
func Inspect2[K, V any](s iter.Seq2[K, V], f func(K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			f(k, v)
			if !yield(k, v) {
				return
			}
		}
	}
}

// Join turns a key/value sequence into an element sequence by applying f to
// each pair. It is the inverse direction of Split.
//
//	Join(Enumerate(Range1(3)), func(k, v int) string {
//	    return fmt.Sprintf("%d=%d", k, v)
//	})  // yields "0=0", "1=1", "2=2"
func Join[E, K, V any](s iter.Seq2[K, V], f func(K, V) E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Keys returns a sequence of the keys from a key/value sequence.
//
//	Keys(Enumerate(seqOf("a", "b")))  // yields 0, 1
func Keys[K, V any](s iter.Seq2[K, V]) iter.Seq[K] {
	return Join(s, func(k K, _ V) K { return k })
}

// Values returns a sequence of the values from a key/value sequence.
//
//	Values(Enumerate(seqOf("a", "b")))  // yields "a", "b"
func Values[K, V any](s iter.Seq2[K, V]) iter.Seq[V] {
	return Join(s, func(_ K, v V) V { return v })
}

// Swap swaps each key/value pair into a (value, key) pair.
//
//	Swap(Enumerate(seqOf("a", "b")))  // yields ("a",0), ("b",1)
func Swap[K, V any](s iter.Seq2[K, V]) iter.Seq2[V, K] {
	return Map2(s, func(k K, v V) (V, K) { return v, k })
}

// ============================================================================
// Filter / Slice
// ============================================================================

// Filter2 yields only key/value pairs that satisfy f. Pairs for which f
// returns false are skipped. Iteration stops as soon as yield returns false.
//
//	Filter2(Enumerate(seqOf("a", "b", "c")), func(k int, _ string) bool { return k%2 == 0 })
//	// yields (0,"a"), (2,"c")
func Filter2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// FilterMap2 applies f to each pair and keeps only the results for which f
// returned ok=true. It is equivalent to Filter2 followed by Map2, but in a
// single pass. Iteration stops as soon as yield returns false.
//
//	FilterMap2(Enumerate(Range1(5)), func(k, v int) (int, int, bool) {
//	    if v%2 == 0 { return k, v * 10, true }
//	    return 0, 0, false
//	})  // yields (0,0), (2,20), (4,40)
func FilterMap2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2, bool)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range s {
			if k2, v2, ok := f(k1, v1); ok && !yield(k2, v2) {
				return
			}
		}
	}
}

// Take2 yields the first n pairs, then stops. When n <= 0 the result is
// empty. When the source has fewer than n pairs, all of them are yielded.
//
//	Take2(Enumerate(Range1(10)), 3)  // yields (0,0), (1,1), (2,2)
func Take2[K, V any](s iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			return
		}
		for k, v := range s {
			if !yield(k, v) {
				return
			}
			n--
			if n <= 0 {
				return
			}
		}
	}
}

// TakeWhile2 yields pairs while f returns true, then stops. The first pair
// for which f returns false is not yielded and ends the sequence.
//
//	TakeWhile2(Enumerate(Range1(10)), func(k, _ int) bool { return k < 3 })
//	// yields (0,0), (1,1), (2,2)
func TakeWhile2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if !f(k, v) || !yield(k, v) {
				return
			}
		}
	}
}

// Skip2 drops the first n pairs and yields the rest. When n <= 0 nothing is
// skipped. When the source has fewer than n pairs the result is empty.
//
//	Skip2(Enumerate(Range1(5)), 2)  // yields (2,2), (3,3), (4,4)
func Skip2[K, V any](s iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if n > 0 {
				n--
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// SkipWhile2 drops pairs while f returns true, then yields the rest
// (including the first pair for which f returned false) unchanged. Once the
// predicate fails it is never consulted again.
//
//	SkipWhile2(Enumerate(Range1(5)), func(k, _ int) bool { return k < 3 })
//	// yields (3,3), (4,4)
func SkipWhile2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		skip := true
		for k, v := range s {
			if skip {
				skip = f(k, v)
				if skip {
					continue
				}
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// Chain2 concatenates seq1 and seq2 into a single key/value sequence: all
// pairs of seq1 first, then all pairs of seq2. When yield returns false,
// iteration of the current source stops immediately and the other source is
// never consumed.
//
//	Chain2(Once2("a", 1), Once2("b", 2))  // yields ("a",1), ("b",2)
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

// ============================================================================
// Terminal
// ============================================================================

// ForEach2 consumes the sequence and calls f for each key/value pair. It is a
// terminal operation with no return value, typically used for side effects.
// The sequence is fully consumed.
func ForEach2[K, V any](s iter.Seq2[K, V], f func(K, V)) {
	for k, v := range s {
		f(k, v)
	}
}

// TryForEach2 consumes the sequence and calls f for each pair until f returns
// an error. It returns the first error and stops consuming the sequence
// immediately. Returns nil when f never errors.
func TryForEach2[K, V any](s iter.Seq2[K, V], f func(K, V) error) error {
	for k, v := range s {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Fold2 reduces the sequence to a single value by repeatedly applying
// f(acc, key, value), starting from init. Returns init unchanged when the
// sequence is empty.
//
//	Fold2(Enumerate(Range1(5)), 0, func(acc, k, v int) int { return acc + v })
//	// returns 10
func Fold2[K, V, A any](s iter.Seq2[K, V], init A, f func(A, K, V) A) A {
	for k, v := range s {
		init = f(init, k, v)
	}
	return init
}

// Reduce2 reduces the sequence to a single key/value pair using the first
// pair as the initial accumulator and f to combine it with each subsequent
// pair. Returns (zero, zero, false) when the sequence is empty.
//
//	Reduce2(Enumerate(Range1(5)), func(k1, v1, k2, v2 int) (int, int) {
//	    return k1 + k2, v1 + v2
//	})  // returns (10, 10, true)
func Reduce2[K, V any](s iter.Seq2[K, V], f func(K, V, K, V) (K, V)) (K, V, bool) {
	it, stop := iter.Pull2(s)
	defer stop()

	accK, accV, ok := it()
	if !ok {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	for k, v, ok := it(); ok; k, v, ok = it() {
		accK, accV = f(accK, accV, k, v)
	}
	return accK, accV, true
}

// TryReduce2 is like Reduce2 but stops when f returns an error. Returns
// (zero, zero, false, nil) when the sequence is empty. On error returns the
// accumulator pair built so far, ok=true, and the error.
func TryReduce2[K, V any](s iter.Seq2[K, V], f func(K, V, K, V) (K, V, error)) (K, V, bool, error) {
	it, stop := iter.Pull2(s)
	defer stop()

	accK, accV, ok := it()
	if !ok {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false, nil
	}
	for k, v, ok := it(); ok; k, v, ok = it() {
		var err error
		accK, accV, err = f(accK, accV, k, v)
		if err != nil {
			return accK, accV, true, err
		}
	}
	return accK, accV, true, nil
}

// TryFold2 is like Fold2 but stops when f returns an error. Returns the
// accumulator built so far together with the first error; on success returns
// (final accumulator, nil).
func TryFold2[K, V, A any](s iter.Seq2[K, V], init A, f func(A, K, V) (A, error)) (A, error) {
	for k, v := range s {
		var err error
		init, err = f(init, k, v)
		if err != nil {
			return init, err
		}
	}
	return init, nil
}

// Size2 counts the number of key/value pairs in the sequence by fully
// consuming it. Beware: calling Size2 on an infinite sequence never returns.
func Size2[K, V any](s iter.Seq2[K, V]) int {
	count := 0
	for range s {
		count++
	}
	return count
}

// SizeFunc2 counts the pairs that satisfy f. The sequence is fully consumed.
//
//	SizeFunc2(Enumerate(Range1(10)), func(k, _ int) bool { return k%2 == 0 })
//	// returns 5
func SizeFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) int {
	count := 0
	for k, v := range s {
		if f(k, v) {
			count++
		}
	}
	return count
}

// SizeValue2 counts pairs whose key equals k and value equals v. It is a
// shortcut for SizeFunc2 with an equality predicate.
//
//	SizeValue2(Enumerate(seqOf("a", "b", "a")), 0, "a")  // returns 1
func SizeValue2[K, V comparable](s iter.Seq2[K, V], k K, v V) int {
	return SizeFunc2(s, func(ck K, cv V) bool { return ck == k && cv == v })
}

// ============================================================================
// Compare / Search
// ============================================================================

// Contains2 reports whether the sequence contains the pair (k, v). Stops as
// soon as a match is found. Returns false for an empty sequence. It is a
// shortcut for ContainsFunc2 with an equality predicate.
//
//	Contains2(Enumerate(seqOf("a", "b")), 1, "b")  // returns true
func Contains2[K, V comparable](s iter.Seq2[K, V], k K, v V) bool {
	return ContainsFunc2(s, func(ck K, cv V) bool { return ck == k && cv == v })
}

// ContainsFunc2 reports whether any pair satisfies f. Stops as soon as f
// returns true. Returns false for an empty sequence.
func ContainsFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return true
		}
	}
	return false
}

// Any2 reports whether at least one pair satisfies f. Stops as soon as f
// returns true. Returns false for an empty sequence. Any2 is semantically
// equivalent to ContainsFunc2.
//
//	Any2(Enumerate(seqOf("a", "b")), func(k int, _ string) bool { return k > 0 })
//	// returns true
func Any2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return true
		}
	}
	return false
}

// All2 reports whether every pair satisfies f. Stops as soon as f returns
// false. Returns true for an empty sequence (vacuous truth).
//
//	All2(Enumerate(seqOf(1, 2, 3)), func(_ int, v int) bool { return v < 5 })
//	// returns true
func All2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if !f(k, v) {
			return false
		}
	}
	return true
}

// First2 returns the first key/value pair of the sequence. If the sequence is
// empty, it returns the zero values and false. Only consumes the first pair.
func First2[K, V any](s iter.Seq2[K, V]) (K, V, bool) {
	for k, v := range s {
		return k, v, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// FirstFunc2 returns the first pair that satisfies f. If no such pair exists,
// it returns the zero values and false. Consumption stops at the first match.
func FirstFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) (K, V, bool) {
	for k, v := range s {
		if f(k, v) {
			return k, v, true
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Last2 returns the last key/value pair of the sequence. If the sequence is
// empty, it returns the zero values and false. The entire sequence is
// consumed.
func Last2[K, V any](s iter.Seq2[K, V]) (K, V, bool) {
	var lastK K
	var lastV V
	found := false
	for k, v := range s {
		lastK = k
		lastV = v
		found = true
	}
	return lastK, lastV, found
}

// LastFunc2 returns the last pair that satisfies f. If no such pair exists,
// it returns the zero values and false. The entire sequence is consumed.
func LastFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) (K, V, bool) {
	var lastK K
	var lastV V
	found := false
	for k, v := range s {
		if f(k, v) {
			lastK = k
			lastV = v
			found = true
		}
	}
	return lastK, lastV, found
}

// Position2 returns the zero-based index of the first pair satisfying f.
// Returns (-1, false) when no pair matches. Consumption stops at the first
// match.
//
//	Position2(Enumerate(seqOf("a", "b", "c")), func(_ int, v string) bool { return v == "b" })
//	// returns (1, true)
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

// Compare2 performs a lexicographic comparison of x and y. Keys are compared
// first, and values are compared only when keys are equal. Returns a negative
// value when x < y, zero when equal, and a positive value when x > y. A
// shorter prefix-equal sequence is "less" than a longer one. Comparison stops
// at the first differing pair or when either sequence is exhausted.
//
//	Compare2(Once2(1, "a"), Once2(1, "b"))  // returns -1 (a < b)
//	Compare2(Once2(1, "a"), Once2(1, "a"))  // returns 0  (equal)
//	Compare2(Once2(1, "a"), Empty2[int, string]())  // returns 1 (x is longer)
func Compare2[K, V cmp.Ordered](x, y iter.Seq2[K, V]) int {
	return CompareFunc2(x, y, func(k1 K, v1 V, k2 K, v2 V) int {
		if c := cmp.Compare(k1, k2); c != 0 {
			return c
		}
		return cmp.Compare(v1, v2)
	})
}

// CompareFunc2 performs a lexicographic comparison of x and y using f on each
// pair of (key, value) tuples. f must follow the cmp.Compare convention:
// negative when a < b, zero when equal, positive when a > b. Returns:
//   - 0 if the sequences are equal (same length and pairs)
//   - a negative value if x is a prefix of y, or the first differing pair of
//     x compares less
//   - a positive value if y is a prefix of x, or the first differing pair of
//     x compares greater
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

// Equal2 reports whether x and y have the same length and elementwise-equal
// keys and values. Stops at the first difference.
func Equal2[K, V comparable](x, y iter.Seq2[K, V]) bool {
	return EqualFunc2(x, y, func(k1 K, v1 V, k2 K, v2 V) bool { return k1 == k2 && v1 == v2 })
}

// EqualFunc2 reports whether x and y have the same length and f(k1, v1, k2,
// v2) holds for every corresponding pair. Stops at the first mismatch.
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
