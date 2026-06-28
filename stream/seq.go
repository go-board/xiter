package stream

import (
	"iter"

	"github.com/go-board/xiter"
)

// Seq is a chainable wrapper around the bare iterator function
// func(yield func(E) bool) (the same underlying type as iter.Seq[E]). It
// exposes the same lazy, single-pass semantics as iter.Seq, but adds methods so
// that pipelines can be written in a fluent style. All non-terminal methods
// return a new Seq that defers all work until iteration; terminal methods
// consume the sequence and produce a value. Every operation honors the
// yield-return-false early stop signal and releases the source as soon as the
// consumer breaks early.
//
// Construct a Seq with Of, FromFunc, or Iterate; chain transform and filter
// methods (Filter, Take, Skip, Inspect, Chain, ...); terminate with a method
// such as ForEach, Reduce, First, or Fold.
//
// Methods that change the element type — Map, MapWhile, FilterMap, Split, Zip,
// ZipWith, Fold, TryFold, Scan — require Go 1.27 method-level generics and are
// only available when building with Go 1.27 or newer.
//
//	s := Of(xiter.Range1(10)).
//	    Filter(func(n int) bool { return n%2 == 0 }).
//	    Map(strconv.Itoa)               // requires Go 1.27 method-level generics
//	    Take(3)
//	s.ForEach(func(s string) { fmt.Println(s) })
type Seq[E any] func(yield func(E) bool)

// Of wraps a bare iterator function of the form func(yield func(E) bool) as a
// chainable Seq. The function is used as-is; no copy is made. The argument has
// the same underlying type as iter.Seq[E], so an iter.Seq can be passed
// directly with an explicit conversion or via Iter().
func Of[E any](s func(yield func(E) bool)) Seq[E] { return Seq[E](s) }

// FromFunc generates a Seq by repeatedly calling f until it returns ok=false.
// f is called lazily, only as the returned Seq is consumed. Useful for turning
// a callback-style or pull-style supplier into a sequence.
//
//	FromFunc(func() (int, bool) { return readLine(reader) })
func FromFunc[E any](f func() (E, bool)) Seq[E] { return Of(xiter.FromFunc(f)) }

// Iterate generates a Seq starting from seed, applying next to derive each
// subsequent element. The sequence ends when next returns ok=false; the value
// returned alongside ok=false is discarded. No external limiter is required.
//
//	Iterate(1, func(x int) (int, bool) {
//	    if x >= 16 { return 0, false }
//	    return x * 2, true
//	})  // Seq[int] yielding 1, 2, 4, 8, 16
func Iterate[E any](seed E, next func(E) (E, bool)) Seq[E] {
	return Of(xiter.Iterate(seed, next))
}

// Iter unwraps the Seq back to a plain iter.Seq so it can be passed to package
// xiter functions or used in a range-over-func loop directly.
func (s Seq[E]) Iter() iter.Seq[E] { return iter.Seq[E](s) }

// Filter returns a Seq that yields only elements for which f returns true.
// Elements are dropped lazily; the source is not consumed until the returned
// Seq is iterated.
func (s Seq[E]) Filter(f func(E) bool) Seq[E] { return Of(xiter.Filter(s.Iter(), f)) }

// Inspect returns a Seq that calls f for each element as it passes through,
// yielding the element unchanged. f is invoked only for elements actually
// visited when the returned Seq is consumed. Useful for debugging or logging
// inside a lazy pipeline without consuming the sequence.
func (s Seq[E]) Inspect(f func(E)) Seq[E] { return Of(xiter.Inspect(s.Iter(), f)) }

// Take returns a Seq yielding the first n elements, then stops. When n <= 0
// the result is empty. When the source has fewer than n elements, all of them
// are yielded. The source is released as soon as n elements have been yielded
// or the consumer breaks early.
func (s Seq[E]) Take(n int) Seq[E] { return Of(xiter.Take(s.Iter(), n)) }

// Skip returns a Seq that drops the first n elements and yields the rest.
// When n <= 0 nothing is skipped. When the source has fewer than n elements
// the result is empty.
func (s Seq[E]) Skip(n int) Seq[E] { return Of(xiter.Skip(s.Iter(), n)) }

// TakeWhile returns a Seq that yields elements while f returns true and stops
// at the first element for which f returns false (that element is not yielded).
// The source is released as soon as f returns false or the consumer breaks
// early.
//
//	Of(seqOf(1, 2, 3, 4, 1)).TakeWhile(func(n int) bool { return n < 3 })
//	// yields 1, 2
func (s Seq[E]) TakeWhile(f func(E) bool) Seq[E] { return Of(xiter.TakeWhile(s.Iter(), f)) }

// SkipWhile returns a Seq that drops elements while f returns true, then
// yields the rest unchanged (including the first element for which f returned
// false). Once the predicate fails it is never consulted again.
//
//	Of(seqOf(1, 2, 3, 4, 1)).SkipWhile(func(n int) bool { return n < 3 })
//	// yields 3, 4, 1
func (s Seq[E]) SkipWhile(f func(E) bool) Seq[E] { return Of(xiter.SkipWhile(s.Iter(), f)) }

// Chain concatenates s and other into a single Seq: all elements of s first,
// then all elements of other. Either side may be empty or infinite.
//
//	Of(seqOf(1, 2)).Chain(Of(seqOf(10, 11)))  // yields 1, 2, 10, 11
func (s Seq[E]) Chain(other Seq[E]) Seq[E] { return Of(xiter.Chain(s.Iter(), other.Iter())) }

// Enumerate converts s into a Seq2[int, E] whose key is the zero-based index
// of each element and whose value is the element itself.
//
//	Of(xiter.Range1(3)).Enumerate()  // yields (0,0), (1,1), (2,2)
func (s Seq[E]) Enumerate() Seq2[int, E] { return Of2(xiter.Enumerate(s.Iter())) }

// ForEach is a terminal operation that consumes s and calls f for each element.
// It has no return value. The sequence is fully consumed unless f panics.
func (s Seq[E]) ForEach(f func(E)) { xiter.ForEach(s.Iter(), f) }

// TryForEach is a terminal operation that consumes s and calls f for each
// element until f returns an error. It returns the first error and stops
// consuming the sequence immediately. Returns nil when f never errors.
func (s Seq[E]) TryForEach(f func(E) error) error { return xiter.TryForEach(s.Iter(), f) }

// Reduce is a terminal operation that reduces s using the first element as the
// initial accumulator and f to combine it with each subsequent element. Returns
// (zero, false) when the sequence is empty. The entire sequence is consumed.
//
//	Of(seqOf(1, 2, 3, 4)).Reduce(func(a, b int) int { return a + b })  // (10, true)
func (s Seq[E]) Reduce(f func(E, E) E) (E, bool) { return xiter.Reduce(s.Iter(), f) }

// TryReduce is like Reduce but stops when f returns an error. Returns
// (zero, false, nil) when the sequence is empty. On error returns the
// accumulator built so far, ok=true, and the error.
func (s Seq[E]) TryReduce(f func(E, E) (E, error)) (E, bool, error) {
	return xiter.TryReduce(s.Iter(), f)
}

// Size is a terminal operation that fully consumes s and returns the number of
// elements. Beware: calling Size on an infinite sequence never returns.
func (s Seq[E]) Size() int { return xiter.Size(s.Iter()) }

// SizeFunc is a terminal operation that counts elements satisfying f. The
// sequence is fully consumed.
func (s Seq[E]) SizeFunc(f func(E) bool) int { return xiter.SizeFunc(s.Iter(), f) }

// Any is a terminal operation that reports whether at least one element
// satisfies f. It stops as soon as f returns true. Returns false for an empty
// sequence.
func (s Seq[E]) Any(f func(E) bool) bool { return xiter.Any(s.Iter(), f) }

// All is a terminal operation that reports whether every element satisfies f.
// It stops as soon as f returns false. Returns true for an empty sequence
// (vacuous truth).
func (s Seq[E]) All(f func(E) bool) bool { return xiter.All(s.Iter(), f) }

// First is a terminal operation that returns the first element. If the
// sequence is empty, it returns the zero value and false. Only the first
// element is consumed.
func (s Seq[E]) First() (E, bool) { return xiter.First(s.Iter()) }

// Last is a terminal operation that returns the last element. If the sequence
// is empty, it returns the zero value and false. The entire sequence is
// consumed.
func (s Seq[E]) Last() (E, bool) { return xiter.Last(s.Iter()) }

// FirstFunc is a terminal operation that returns the first element satisfying
// f. If no such element exists, it returns the zero value and false.
// Consumption stops at the first match.
func (s Seq[E]) FirstFunc(f func(E) bool) (E, bool) { return xiter.FirstFunc(s.Iter(), f) }

// LastFunc is a terminal operation that returns the last element satisfying f.
// If no such element exists, it returns the zero value and false. The entire
// sequence is consumed.
func (s Seq[E]) LastFunc(f func(E) bool) (E, bool) { return xiter.LastFunc(s.Iter(), f) }

// Position is a terminal operation that returns the zero-based index of the
// first element satisfying f. Returns (-1, false) when no element matches.
// Consumption stops at the first match.
func (s Seq[E]) Position(f func(E) bool) (int, bool) { return xiter.Position(s.Iter(), f) }

// IsSortedFunc is a terminal operation that reports whether s is sorted by
// comparator f. The comparison accepts both non-decreasing and non-increasing
// order; the sequence is considered sorted when every adjacent pair is
// consistently ordered in one of those directions.
//
//	Of(seqOf(1, 2, 3)).IsSortedFunc(cmp.Compare)  // true  (non-decreasing)
//	Of(seqOf(3, 2, 1)).IsSortedFunc(cmp.Compare)  // true  (non-increasing)
//	Of(seqOf(1, 3, 2)).IsSortedFunc(cmp.Compare)  // false (inconsistent)
func (s Seq[E]) IsSortedFunc(f func(E, E) int) bool { return xiter.IsSortedFunc(s.Iter(), f) }

// CompareFunc is a terminal operation that lexicographically compares s and
// other element by element using f, following cmp.Compare's negative/zero/
// positive convention. Comparison stops at the first differing element; if one
// sequence is a prefix of the other, the shorter sequence is considered
// smaller. Returns 0 when the two sequences are equal.
func (s Seq[E]) CompareFunc(other Seq[E], f func(E, E) int) int {
	return xiter.CompareFunc(s.Iter(), other.Iter(), f)
}

// EqualFunc is a terminal operation that reports whether s and other have the
// same length and every pair of corresponding elements satisfies f. Both
// sequences are fully consumed.
func (s Seq[E]) EqualFunc(other Seq[E], f func(E, E) bool) bool {
	return xiter.EqualFunc(s.Iter(), other.Iter(), f)
}

// MaxFunc is a terminal operation that returns the maximum element of s by
// comparator f. Returns (zero, false) when the sequence is empty. The entire
// sequence is consumed.
func (s Seq[E]) MaxFunc(f func(E, E) int) (E, bool) { return xiter.MaxFunc(s.Iter(), f) }

// MinFunc is a terminal operation that returns the minimum element of s by
// comparator f. Returns (zero, false) when the sequence is empty. The entire
// sequence is consumed.
func (s Seq[E]) MinFunc(f func(E, E) int) (E, bool) { return xiter.MinFunc(s.Iter(), f) }

// MinMaxFunc is a terminal operation that returns both the minimum and the
// maximum element of s by comparator f in a single pass. Returns
// (zero, zero, false) when the sequence is empty. The entire sequence is
// consumed.
//
//	Of(seqOf(3, 1, 4, 1, 5, 9)).MinMaxFunc(cmp.Compare)  // (1, 9, true)
func (s Seq[E]) MinMaxFunc(f func(E, E) int) (E, E, bool) {
	return xiter.MinMaxFunc(s.Iter(), f)
}

// ContainsFunc is a terminal operation that reports whether any element of s
// satisfies f. It stops as soon as f returns true. Returns false for an empty
// sequence. Equivalent to Any but named to mirror the Contains family.
func (s Seq[E]) ContainsFunc(f func(E) bool) bool { return xiter.ContainsFunc(s.Iter(), f) }
