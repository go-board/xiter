package xiter

import (
	"cmp"
	"iter"
)

// ============================================================================
// Source
// ============================================================================

// Range1 generates an integer sequence from 0 to end-1 (end not included).
// Returns an empty sequence when end <= 0.
//
//	Range1(5)  // yields 0, 1, 2, 3, 4
//	Range1(0)  // yields nothing
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
// Returns an empty sequence when start >= end.
//
//	Range2(2, 5)  // yields 2, 3, 4
//	Range2(5, 2)  // yields nothing (start >= end)
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
// (end not included). A positive step counts upward; a negative step counts
// downward. Returns an empty sequence when step is 0, or when the step
// direction is inconsistent with the start/end relationship.
//
//	Range3(1, 10, 2)   // yields 1, 3, 5, 7, 9
//	Range3(10, 1, -2)  // yields 10, 8, 6, 4, 2
//	Range3(1, 10, 0)   // yields nothing (step == 0)
//	Range3(10, 1, 2)   // yields nothing (positive step, start >= end)
func Range3[N integral](start, end, step N) iter.Seq[N] {
	if step == 0 {
		return Empty[N]()
	}
	if step > 0 && start >= end {
		return Empty[N]()
	}
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

// FromFunc generates a sequence from a supplier function that returns
// (element, continue). The sequence ends when continue is false. The supplier
// is called lazily, only as the returned sequence is consumed.
//
//	i := 0
//	FromFunc(func() (int, bool) {
//	    i++
//	    if i > 3 { return 0, false }
//	    return i, true
//	})  // yields 1, 2, 3
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

// Iterate generates a sequence starting from seed, where each subsequent
// element is produced by applying next to the previous one. The sequence ends
// when next returns ok=false; the value returned alongside ok=false is
// discarded. It is the Go analog of Java's Stream.iterate, with the
// termination condition folded into next so no external limiter is required.
//
//	Iterate(1, func(x int) (int, bool) {
//	    if x >= 16 { return 0, false }
//	    return x * 2, true
//	})  // yields 1, 2, 4, 8, 16
func Iterate[E any](seed E, next func(E) (E, bool)) iter.Seq[E] {
	return func(yield func(E) bool) {
		e := seed
		if !yield(e) {
			return
		}
		for {
			ne, ok := next(e)
			if !ok {
				return
			}
			if !yield(ne) {
				return
			}
			e = ne
		}
	}
}

// Once generates a sequence containing a single element. The returned
// sequence yields exactly one value and then ends; yield's return value is
// ignored since there is nothing further to release.
func Once[E any](e E) iter.Seq[E] {
	return func(yield func(E) bool) {
		yield(e)
	}
}

// Empty generates an empty sequence. Iterating it yields nothing and returns
// immediately.
func Empty[E any]() iter.Seq[E] {
	return func(yield func(E) bool) {
	}
}

// Repeat generates an infinite sequence repeating a single element. Always
// pair it with a limiting operator such as Take to avoid consuming forever.
//
//	Take(Repeat(1), 3)  // yields 1, 1, 1
func Repeat[E any](e E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			if !yield(e) {
				return
			}
		}
	}
}

// ============================================================================
// Transform
// ============================================================================

// Map applies f to each element and yields the results. The output length
// matches the input length. Evaluation is lazy: f is invoked only as the
// returned sequence is consumed, and iteration stops as soon as yield returns
// false.
//
//	Map(Range1(3), func(x int) int { return x * 2 })  // yields 0, 2, 4
func Map[E1, E2 any](s iter.Seq[E1], f func(E1) E2) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e := range s {
			if !yield(f(e)) {
				return
			}
		}
	}
}

// MapWhile applies f to each element and yields the results until f returns
// ok=false, at which point the sequence stops. The element that caused the
// stop is not yielded.
//
//	MapWhile(Range1(10), func(x int) (int, bool) {
//	    if x < 3 { return x * 10, true }
//	    return 0, false
//	})  // yields 0, 10, 20
func MapWhile[E1, E2 any](s iter.Seq[E1], f func(E1) (E2, bool)) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e1 := range s {
			if e2, ok := f(e1); !ok || !yield(e2) {
				return
			}
		}
	}
}

// FlatMap applies f to each element — which itself returns a sequence — and
// flattens the results into a single sequence. The total length is the sum of
// the lengths of all inner sequences. When yield returns false, both the
// current inner sequence and the outer source are released immediately.
//
//	FlatMap(Range1(3), func(x int) iter.Seq[int] {
//	    return Range1(x + 1)
//	})  // yields 0, 0,1, 0,1,2
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

// Flatten flattens a sequence of sequences into a single sequence. It is a
// shortcut for FlatMap with the identity function, and shares the same
// early-termination behavior.
//
//	Flatten(seqOf(seqOf(1, 2), seqOf(3, 4)))  // yields 1, 2, 3, 4
func Flatten[E any](s iter.Seq[iter.Seq[E]]) iter.Seq[E] {
	return FlatMap(s, func(e iter.Seq[E]) iter.Seq[E] { return e })
}

// Inspect calls f for each element as it passes through and yields the element
// unchanged. Useful for debugging or logging in the middle of a lazy pipeline
// without affecting the data. f is called only when the returned sequence is
// consumed, and only for elements actually visited.
func Inspect[E any](s iter.Seq[E], f func(E)) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			f(e)
			if !yield(e) {
				return
			}
		}
	}
}

// Enumerate pairs each element with its zero-based index, yielding (index,
// element) pairs. The index starts at 0 for the first element.
//
//	Enumerate(seqOf("a", "b"))  // yields (0,"a"), (1,"b")
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

// Split turns an element sequence into a key/value sequence by applying f to
// each element. It is the inverse direction of Join.
//
//	Split(seqOf(1, 2), func(e int) (int, int) { return e, e*10 })
//	// yields (1,10), (2,20)
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

// Cast attempts a type assertion of each any element to E, yielding (value, ok)
// pairs. ok is true when the assertion succeeded; when it failed the zero
// value of E is yielded alongside false. No elements are dropped.
//
//	Cast[int](seqOf[any](1, "x", 3))
//	// yields (1,true), (0,false), (3,true)
func Cast[E any](s iter.Seq[any]) iter.Seq2[E, bool] {
	return func(yield func(E, bool) bool) {
		for a := range s {
			if e, ok := a.(E); !yield(e, ok) {
				return
			}
		}
	}
}

// Scan maintains an accumulator seeded with init and applies f to it together
// with each element. The updated accumulator is yielded for every element, so
// the output length matches the input length. When f returns ok=false, the
// scan stops early and no value is yielded for the current element.
//
// Unlike Fold, Scan exposes every intermediate accumulator state instead of
// only the final result. Unlike MapWhile, Scan carries mutable state across
// elements.
//
//	Scan(Range1(5), 0, func(acc, e int) (int, bool) {
//	    return acc + e, true
//	})  // yields 0, 1, 3, 6, 10  (running sum)
func Scan[E, A any](s iter.Seq[E], init A, f func(A, E) (A, bool)) iter.Seq[A] {
	return func(yield func(A) bool) {
		acc := init
		for e := range s {
			next, ok := f(acc, e)
			if !ok {
				return
			}
			acc = next
			if !yield(acc) {
				return
			}
		}
	}
}

// ============================================================================
// Filter / Slice
// ============================================================================

// Filter yields only elements that satisfy f. Elements for which f returns
// false are skipped. Evaluation is lazy: f is invoked only as the returned
// sequence is consumed, and iteration stops as soon as yield returns false.
//
//	Filter(Range1(5), func(x int) bool { return x%2 == 0 })  // yields 0, 2, 4
func Filter[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if f(e) && !yield(e) {
				return
			}
		}
	}
}

// FilterMap applies f to each element and keeps only the results for which f
// returned ok=true. It is equivalent to Filter followed by Map, but in a
// single pass. Iteration stops as soon as yield returns false.
//
//	FilterMap(Range1(5), func(x int) (int, bool) {
//	    if x%2 == 0 { return x * 2, true }
//	    return 0, false
//	})  // yields 0, 4, 8
func FilterMap[E1, E2 any](s iter.Seq[E1], f func(E1) (E2, bool)) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e1 := range s {
			if e2, ok := f(e1); ok && !yield(e2) {
				return
			}
		}
	}
}

// Take yields the first n elements, then stops. When n <= 0 the result is
// empty. When the source has fewer than n elements, all of them are yielded.
//
//	Take(Range1(10), 3)  // yields 0, 1, 2
//	Take(Range1(2), 5)   // yields 0, 1
//	Take(Range1(5), 0)   // yields nothing
func Take[E any](s iter.Seq[E], n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		if n <= 0 {
			return
		}
		for e := range s {
			if !yield(e) {
				return
			}
			n--
			if n <= 0 {
				return
			}
		}
	}
}

// TakeWhile yields elements while f returns true, then stops. The first
// element for which f returns false is not yielded and ends the sequence.
//
//	TakeWhile(Range1(10), func(x int) bool { return x < 3 })  // yields 0, 1, 2
func TakeWhile[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if !f(e) || !yield(e) {
				return
			}
		}
	}
}

// Skip drops the first n elements and yields the rest. When n <= 0 nothing is
// skipped. When the source has fewer than n elements the result is empty.
//
//	Skip(Range1(5), 2)  // yields 2, 3, 4
//	Skip(Range1(3), 10) // yields nothing
func Skip[E any](s iter.Seq[E], n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if n > 0 {
				n--
				continue
			}
			if !yield(e) {
				return
			}
		}
	}
}

// SkipWhile drops elements while f returns true, then yields the rest
// (including the first element for which f returned false) unchanged. Once
// the predicate fails it is never consulted again.
//
//	SkipWhile(Range1(5), func(x int) bool { return x < 3 })
//	// yields 3, 4
func SkipWhile[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		skip := true
		for e := range s {
			if skip {
				skip = f(e)
				if skip {
					continue
				}
			}
			if !yield(e) {
				return
			}
		}
	}
}

// Chain concatenates seq1 and seq2 into a single sequence: all elements of
// seq1 first, then all elements of seq2. When yield returns false, iteration
// of the current source stops immediately and the other source is never
// consumed.
//
//	Chain(Range1(2), Range2(10, 12))  // yields 0, 1, 10, 11
func Chain[E any](seq1, seq2 iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range seq1 {
			if !yield(e) {
				return
			}
		}
		for e := range seq2 {
			if !yield(e) {
				return
			}
		}
	}
}

// Zip pairs elements from x and y position-by-position, yielding (x_i, y_i)
// pairs. Iteration stops as soon as either sequence is exhausted; excess
// elements of the longer sequence are never consumed.
//
//	Zip(Range1(5), Range2(10, 13))
//	// yields (0,10), (1,11), (2,12)
func Zip[E1, E2 any](x iter.Seq[E1], y iter.Seq[E2]) iter.Seq2[E1, E2] {
	return func(yield func(E1, E2) bool) {
		nextX, stopX := iter.Pull(x)
		defer stopX()
		nextY, stopY := iter.Pull(y)
		defer stopY()

		for {
			e1, ok1 := nextX()
			if !ok1 {
				return
			}
			e2, ok2 := nextY()
			if !ok2 || !yield(e1, e2) {
				return
			}
		}
	}
}

// ZipWith combines elements from x and y position-by-position via f, yielding
// f(x_i, y_i). Iteration stops as soon as either sequence is exhausted.
//
//	ZipWith(Range1(3), Range2(10, 13), func(a, b int) int { return a + b })
//	// yields 10, 12, 14
func ZipWith[E1, E2, E3 any](x iter.Seq[E1], y iter.Seq[E2], f func(E1, E2) E3) iter.Seq[E3] {
	return func(yield func(E3) bool) {
		for e1, e2 := range Zip(x, y) {
			if !yield(f(e1, e2)) {
				return
			}
		}
	}
}

// ============================================================================
// Terminal
// ============================================================================

// ForEach consumes the sequence and calls f for each element. It is a
// terminal operation with no return value, typically used for side effects.
// The sequence is fully consumed.
func ForEach[E any](s iter.Seq[E], f func(E)) {
	for e := range s {
		f(e)
	}
}

// TryForEach consumes the sequence and calls f for each element until f
// returns an error. It returns the first error and stops consuming the
// sequence immediately. Returns nil when f never errors.
func TryForEach[E any](s iter.Seq[E], f func(E) error) error {
	for e := range s {
		if err := f(e); err != nil {
			return err
		}
	}
	return nil
}

// Fold reduces the sequence to a single value by repeatedly applying
// f(acc, element), starting from init. Returns init unchanged when the
// sequence is empty. Unlike Reduce, Fold accepts an explicit initial
// accumulator and works on empty sequences.
//
//	Fold(Range1(5), 0, func(acc, e int) int { return acc + e })  // returns 10
func Fold[E any, A any](s iter.Seq[E], init A, f func(A, E) A) A {
	for e := range s {
		init = f(init, e)
	}
	return init
}

// Reduce reduces the sequence to a single value using the first element as
// the initial accumulator and f to combine it with each subsequent element.
// Returns (zero, false) when the sequence is empty.
//
//	Reduce(Range1(5), func(a, b int) int { return a + b })  // returns (10, true)
func Reduce[E any](s iter.Seq[E], f func(E, E) E) (E, bool) {
	it, stop := iter.Pull(s)
	defer stop()

	acc, ok := it()
	if !ok {
		var zero E
		return zero, false
	}
	for e, ok := it(); ok; e, ok = it() {
		acc = f(acc, e)
	}
	return acc, true
}

// TryReduce is like Reduce but stops when f returns an error. Returns
// (zero, false, nil) when the sequence is empty. On error returns the
// accumulator built so far, ok=true, and the error.
func TryReduce[E any](s iter.Seq[E], f func(E, E) (E, error)) (E, bool, error) {
	it, stop := iter.Pull(s)
	defer stop()

	acc, ok := it()
	if !ok {
		var zero E
		return zero, false, nil
	}
	for e, ok := it(); ok; e, ok = it() {
		var err error
		acc, err = f(acc, e)
		if err != nil {
			return acc, true, err
		}
	}
	return acc, true, nil
}

// TryFold is like Fold but stops when f returns an error. Returns the
// accumulator built so far together with the first error; on success returns
// (final accumulator, nil).
func TryFold[E any, A any](s iter.Seq[E], init A, f func(A, E) (A, error)) (A, error) {
	for e := range s {
		var err error
		init, err = f(init, e)
		if err != nil {
			return init, err
		}
	}
	return init, nil
}

// Size counts the number of elements in the sequence by fully consuming it.
// Beware: calling Size on an infinite sequence never returns.
func Size[E any](s iter.Seq[E]) int {
	count := 0
	for range s {
		count++
	}
	return count
}

// SizeFunc counts the elements that satisfy f. The sequence is fully
// consumed.
//
//	SizeFunc(Range1(10), func(x int) bool { return x%2 == 0 })  // returns 5
func SizeFunc[E any](s iter.Seq[E], f func(E) bool) int {
	count := 0
	for e := range s {
		if f(e) {
			count++
		}
	}
	return count
}

// SizeValue counts elements equal to v. It is a shortcut for
// SizeFunc(s, func(e E) bool { return e == v }).
//
//	SizeValue(seqOf(1, 2, 2, 3, 2), 2)  // returns 3
func SizeValue[E comparable](s iter.Seq[E], v E) int {
	return SizeFunc(s, func(e E) bool { return e == v })
}

// ============================================================================
// Compare / Search
// ============================================================================

// Contains reports whether the sequence contains v. Stops as soon as a match
// is found. Returns false for an empty sequence. It is a shortcut for
// ContainsFunc with an equality predicate.
//
//	Contains(seqOf(1, 2, 3), 2)  // returns true
//	Contains(seqOf(1, 2, 3), 5)  // returns false
func Contains[E comparable](s iter.Seq[E], v E) bool {
	return ContainsFunc(s, func(e E) bool { return e == v })
}

// ContainsFunc reports whether any element satisfies f. Stops as soon as f
// returns true. Returns false for an empty sequence.
func ContainsFunc[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

// Any reports whether at least one element satisfies f. Stops as soon as f
// returns true. Returns false for an empty sequence. Any is semantically
// equivalent to ContainsFunc.
//
//	Any(seqOf(1, 2, 3), func(x int) bool { return x > 2 })  // returns true
//	Any(Empty[int](), func(x int) bool { return x > 2 })    // returns false
func Any[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

// All reports whether every element satisfies f. Stops as soon as f returns
// false. Returns true for an empty sequence (vacuous truth).
//
//	All(seqOf(1, 2, 3), func(x int) bool { return x < 5 })  // returns true
//	All(seqOf(1, 2, 5), func(x int) bool { return x < 5 })  // returns false
//	All(Empty[int](), func(x int) bool { return x < 5 })    // returns true (vacuous)
func All[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if !f(e) {
			return false
		}
	}
	return true
}

// First returns the first element of the sequence. If the sequence is empty,
// it returns the zero value and false. Only consumes the first element.
func First[E any](s iter.Seq[E]) (E, bool) {
	for e := range s {
		return e, true
	}
	var zero E
	return zero, false
}

// FirstFunc returns the first element that satisfies f. If no such element
// exists, it returns the zero value and false. Consumption stops at the
// first match.
func FirstFunc[E any](s iter.Seq[E], f func(E) bool) (E, bool) {
	for e := range s {
		if f(e) {
			return e, true
		}
	}
	var zero E
	return zero, false
}

// Last returns the last element of the sequence. If the sequence is empty,
// it returns the zero value and false. The entire sequence is consumed.
func Last[E any](s iter.Seq[E]) (E, bool) {
	var last E
	found := false
	for e := range s {
		last = e
		found = true
	}
	return last, found
}

// LastFunc returns the last element that satisfies f. If no such element
// exists, it returns the zero value and false. The entire sequence is
// consumed.
func LastFunc[E any](s iter.Seq[E], f func(E) bool) (E, bool) {
	var last E
	found := false
	for e := range s {
		if f(e) {
			last = e
			found = true
		}
	}
	return last, found
}

// Position returns the zero-based index of the first element satisfying f.
// Returns (-1, false) when no element matches. Consumption stops at the
// first match.
//
//	Position(seqOf("a", "b", "c"), func(s string) bool { return s == "b" })  // returns (1, true)
//	Position(seqOf("a", "b"), func(s string) bool { return s == "z" })       // returns (-1, false)
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

// Compare performs a lexicographic comparison of x and y using cmp.Compare on
// each pair of elements. Returns a negative value when x < y, zero when equal,
// and a positive value when x > y. A shorter prefix-equal sequence is "less"
// than a longer one. Comparison stops at the first differing element or when
// either sequence is exhausted.
//
//	Compare(seqOf(1, 2), seqOf(1, 3))  // returns -1 (2 < 3)
//	Compare(seqOf(1, 2), seqOf(1, 2))  // returns 0  (equal)
//	Compare(seqOf(1, 2), seqOf(1))     // returns 1  (x is longer)
func Compare[E cmp.Ordered](x, y iter.Seq[E]) int {
	return CompareFunc(x, y, cmp.Compare[E])
}

// CompareFunc performs a lexicographic comparison of x and y using f on each
// pair of elements. f must follow the cmp.Compare convention: negative when
// a < b, zero when equal, positive when a > b. Returns:
//   - 0 if the sequences are equal (same length and values)
//   - a negative value if x is a prefix of y, or the first differing element
//     of x compares less
//   - a positive value if y is a prefix of x, or the first differing element
//     of x compares greater
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

// Equal reports whether x and y have the same length and elementwise-equal
// values. Stops at the first difference.
//
//	Equal(seqOf(1, 2, 3), seqOf(1, 2, 3))  // returns true
//	Equal(seqOf(1, 2), seqOf(1, 2, 3))     // returns false (different length)
func Equal[E comparable](x, y iter.Seq[E]) bool {
	return EqualFunc(x, y, func(e1, e2 E) bool { return e1 == e2 })
}

// EqualFunc reports whether x and y have the same length and f(e1, e2) holds
// for every corresponding pair. Stops at the first mismatch.
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

// Max returns the maximum element using cmp.Compare. Returns (zero, false)
// when the sequence is empty. The entire sequence is consumed.
//
//	Max(seqOf(3, 1, 4, 1, 5))  // returns (5, true)
//	Max(Empty[int]())          // returns (0, false)
func Max[E cmp.Ordered](s iter.Seq[E]) (E, bool) {
	return MaxFunc(s, cmp.Compare[E])
}

// MaxFunc returns the maximum element according to cmp, which must follow the
// cmp.Compare convention (negative when a < b, zero when equal, positive
// when a > b). Returns (zero, false) when the sequence is empty.
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

// Min returns the minimum element using cmp.Compare. Returns (zero, false)
// when the sequence is empty. The entire sequence is consumed.
//
//	Min(seqOf(3, 1, 4, 1, 5))  // returns (1, true)
//	Min(Empty[int]())          // returns (0, false)
func Min[E cmp.Ordered](s iter.Seq[E]) (E, bool) {
	return MinFunc(s, cmp.Compare[E])
}

// MinFunc returns the minimum element according to cmp, which must follow the
// cmp.Compare convention (negative when a < b, zero when equal, positive
// when a > b). Returns (zero, false) when the sequence is empty.
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

// MinMax returns both the minimum and maximum element using cmp.Compare, in a
// single pass. Returns (zero, zero, false) when the sequence is empty.
//
//	MinMax(seqOf(3, 1, 4, 1, 5))  // returns (1, 5, true)
//	MinMax(Empty[int]())          // returns (0, 0, false)
func MinMax[E cmp.Ordered](s iter.Seq[E]) (E, E, bool) {
	return MinMaxFunc(s, cmp.Compare[E])
}

// MinMaxFunc returns both the minimum and maximum element according to cmp in
// a single pass. cmp must follow the cmp.Compare convention (negative when
// a < b, zero when equal, positive when a > b). Returns (zero, zero, false)
// when the sequence is empty. More efficient than calling Min and Max
// separately since it walks the sequence once.
func MinMaxFunc[E any](s iter.Seq[E], cmp func(E, E) int) (E, E, bool) {
	it, stop := iter.Pull(s)
	defer stop()
	first, ok := it()
	if !ok {
		var zero E
		return zero, zero, false
	}
	minVal, maxVal := first, first
	for elem, ok := it(); ok; elem, ok = it() {
		if cmp(elem, minVal) < 0 {
			minVal = elem
		}
		if cmp(elem, maxVal) > 0 {
			maxVal = elem
		}
	}
	return minVal, maxVal, true
}

// IsSorted reports whether the elements are in non-decreasing order according
// to cmp.Compare. Returns true for an empty or single-element sequence.
//
//	IsSorted(seqOf(1, 2, 3))  // returns true
//	IsSorted(seqOf(1, 3, 2))  // returns false
func IsSorted[E cmp.Ordered](s iter.Seq[E]) bool {
	return IsSortedFunc(s, cmp.Compare[E])
}

// IsSortedFunc reports whether the sequence is sorted according to a
// consistent order using f. A sequence is considered sorted when every
// adjacent pair compares with the same sign (all non-negative, or all
// non-positive) — i.e. either non-decreasing or non-increasing. Returns true
// for an empty or single-element sequence.
//
//	IsSortedFunc(seqOf(1, 2, 3), cmp.Compare)  // true  (non-decreasing)
//	IsSortedFunc(seqOf(3, 2, 1), cmp.Compare)  // true  (non-increasing)
//	IsSortedFunc(seqOf(1, 3, 2), cmp.Compare)  // false (inconsistent)
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
