package stream

import (
	"iter"

	"github.com/go-board/xiter"
)

// Seq2 is a chainable wrapper around the bare iterator function
// func(yield func(K, V) bool) (the same underlying type as iter.Seq2[K, V]).
// It mirrors Seq but for key/value sequences and exposes the same lazy,
// single-pass semantics. All non-terminal methods return a new Seq2 (or Seq,
// for Keys, Values, and Join) that defers all work until iteration; terminal
// methods consume the sequence and produce a value. Every operation honors the
// yield-return-false early stop signal and releases the source as soon as the
// consumer breaks early.
//
// Construct a Seq2 with Of2, FromFunc2, or Iterate2, or by calling Enumerate on
// a Seq; chain transform and filter methods; terminate with a method such as
// ForEach, Reduce, First, or Fold.
//
// Methods that change the key or value type — Map, MapWhile, FilterMap, Join,
// Fold, TryFold — require Go 1.27 method-level generics and are only available
// when building with Go 1.27 or newer.
//
//	Of(xiter.Range1(3)).
//	    Enumerate().
//	    Filter(func(i, _ int) bool { return i%2 == 0 }).
//	    ForEach(func(i, n int) { fmt.Println(i, n) })
type Seq2[K, V any] func(yield func(K, V) bool)

// Of2 wraps a bare iterator function of the form func(yield func(K, V) bool)
// as a chainable Seq2. The function is used as-is; no copy is made. The
// argument has the same underlying type as iter.Seq2[K, V], so an iter.Seq2
// can be passed directly with an explicit conversion or via Iter().
func Of2[K, V any](s func(yield func(K, V) bool)) Seq2[K, V] {
	return Seq2[K, V](s)
}

// FromFunc2 generates a Seq2 by repeatedly calling f until it returns
// ok=false. f is called lazily, only as the returned Seq2 is consumed. Useful
// for turning a pull-style supplier into a key/value sequence.
func FromFunc2[K, V any](f func() (K, V, bool)) Seq2[K, V] {
	return Of2(xiter.FromFunc2(f))
}

// Iterate2 generates a Seq2 starting from (seedK, seedV), applying next to
// derive each subsequent key/value pair. The sequence ends when next returns
// ok=false; the pair returned alongside ok=false is discarded. No external
// limiter is required.
//
//	Iterate2(0, 1, func(k, v int) (int, int, bool) {
//	    if k >= 3 { return 0, 0, false }
//	    return k + 1, v * 2, true
//	})  // Seq2[int,int] yielding (0,1), (1,2), (2,4), (3,8)
func Iterate2[K, V any](seedK K, seedV V, next func(K, V) (K, V, bool)) Seq2[K, V] {
	return Of2(xiter.Iterate2(seedK, seedV, next))
}

// Iter unwraps the Seq2 back to a plain iter.Seq2 so it can be passed to
// package xiter functions or used in a range-over-func loop directly.
func (s Seq2[K, V]) Iter() iter.Seq2[K, V] { return iter.Seq2[K, V](s) }

// Filter returns a Seq2 that yields only pairs for which f returns true.
// Pairs are dropped lazily; the source is not consumed until the returned Seq2
// is iterated.
func (s Seq2[K, V]) Filter(f func(K, V) bool) Seq2[K, V] {
	return Of2(xiter.Filter2(s.Iter(), f))
}

// Keys returns a Seq of the keys from s, in iteration order.
func (s Seq2[K, V]) Keys() Seq[K] { return Of(xiter.Keys(s.Iter())) }

// Values returns a Seq of the values from s, in iteration order.
func (s Seq2[K, V]) Values() Seq[V] { return Of(xiter.Values(s.Iter())) }

// Swap returns a Seq2[V, K] where each pair's key and value are exchanged.
// Useful for inverting a mapping or preparing for a value-indexed lookup.
func (s Seq2[K, V]) Swap() Seq2[V, K] { return Of2(xiter.Swap(s.Iter())) }

// Inspect returns a Seq2 that calls f for each pair as it passes through,
// yielding the pair unchanged. f is invoked only for pairs actually visited
// when the returned Seq2 is consumed. Useful for debugging or logging inside a
// lazy pipeline.
func (s Seq2[K, V]) Inspect(f func(K, V)) Seq2[K, V] {
	return Of2(xiter.Inspect2(s.Iter(), f))
}

// Take returns a Seq2 yielding the first n pairs, then stops. When n <= 0 the
// result is empty. When the source has fewer than n pairs, all of them are
// yielded.
func (s Seq2[K, V]) Take(n int) Seq2[K, V] { return Of2(xiter.Take2(s.Iter(), n)) }

// Skip returns a Seq2 that drops the first n pairs and yields the rest. When
// n <= 0 nothing is skipped. When the source has fewer than n pairs the result
// is empty.
func (s Seq2[K, V]) Skip(n int) Seq2[K, V] { return Of2(xiter.Skip2(s.Iter(), n)) }

// TakeWhile returns a Seq2 that yields pairs while f returns true and stops at
// the first pair for which f returns false (that pair is not yielded). The
// source is released as soon as f returns false or the consumer breaks early.
func (s Seq2[K, V]) TakeWhile(f func(K, V) bool) Seq2[K, V] {
	return Of2(xiter.TakeWhile2(s.Iter(), f))
}

// SkipWhile returns a Seq2 that drops pairs while f returns true, then yields
// the rest unchanged (including the first pair for which f returned false).
// Once the predicate fails it is never consulted again.
func (s Seq2[K, V]) SkipWhile(f func(K, V) bool) Seq2[K, V] {
	return Of2(xiter.SkipWhile2(s.Iter(), f))
}

// StepBy returns a Seq2 that yields every n-th pair starting from the first
// (index 0). When n <= 0 the result is empty.
//
//	Of2(xiter.Enumerate(xiter.Range1(10))).StepBy(3)
//	// yields (0,0), (3,3), (6,6), (9,9)
func (s Seq2[K, V]) StepBy(n int) Seq2[K, V] { return Of2(xiter.StepBy2(s.Iter(), n)) }

// Chain concatenates s and other into a single Seq2: all pairs of s first,
// then all pairs of other. Either side may be empty or infinite.
func (s Seq2[K, V]) Chain(other Seq2[K, V]) Seq2[K, V] {
	return Of2(xiter.Chain2(s.Iter(), other.Iter()))
}

// ForEach is a terminal operation that consumes s and calls f for each pair.
// It has no return value. The sequence is fully consumed unless f panics.
func (s Seq2[K, V]) ForEach(f func(K, V)) { xiter.ForEach2(s.Iter(), f) }

// TryForEach is a terminal operation that consumes s and calls f for each pair
// until f returns an error. It returns the first error and stops consuming the
// sequence immediately. Returns nil when f never errors.
func (s Seq2[K, V]) TryForEach(f func(K, V) error) error {
	return xiter.TryForEach2(s.Iter(), f)
}

// Reduce is a terminal operation that reduces s using the first pair as the
// initial accumulator and f to combine it with each subsequent pair. Returns
// (zero, zero, false) when the sequence is empty.
func (s Seq2[K, V]) Reduce(f func(K, V, K, V) (K, V)) (K, V, bool) {
	return xiter.Reduce2(s.Iter(), f)
}

// TryReduce is like Reduce but stops when f returns an error. Returns
// (zero, zero, false, nil) when the sequence is empty. On error returns the
// accumulator pair built so far, ok=true, and the error.
func (s Seq2[K, V]) TryReduce(f func(K, V, K, V) (K, V, error)) (K, V, bool, error) {
	return xiter.TryReduce2(s.Iter(), f)
}

// Size is a terminal operation that fully consumes s and returns the number of
// pairs. Beware: calling Size on an infinite sequence never returns.
func (s Seq2[K, V]) Size() int { return xiter.Size2(s.Iter()) }

// SizeFunc is a terminal operation that counts pairs satisfying f. The
// sequence is fully consumed.
func (s Seq2[K, V]) SizeFunc(f func(K, V) bool) int { return xiter.SizeFunc2(s.Iter(), f) }

// Any is a terminal operation that reports whether at least one pair satisfies
// f. It stops as soon as f returns true. Returns false for an empty sequence.
func (s Seq2[K, V]) Any(f func(K, V) bool) bool { return xiter.Any2(s.Iter(), f) }

// All is a terminal operation that reports whether every pair satisfies f. It
// stops as soon as f returns false. Returns true for an empty sequence
// (vacuous truth).
func (s Seq2[K, V]) All(f func(K, V) bool) bool { return xiter.All2(s.Iter(), f) }

// First is a terminal operation that returns the first pair. If the sequence
// is empty, it returns the zero values and false. Only the first pair is
// consumed.
func (s Seq2[K, V]) First() (K, V, bool) { return xiter.First2(s.Iter()) }

// Last is a terminal operation that returns the last pair. If the sequence is
// empty, it returns the zero values and false. The entire sequence is
// consumed.
func (s Seq2[K, V]) Last() (K, V, bool) { return xiter.Last2(s.Iter()) }

// FirstFunc is a terminal operation that returns the first pair satisfying f.
// If no such pair exists, it returns the zero values and false. Consumption
// stops at the first match.
func (s Seq2[K, V]) FirstFunc(f func(K, V) bool) (K, V, bool) {
	return xiter.FirstFunc2(s.Iter(), f)
}

// LastFunc is a terminal operation that returns the last pair satisfying f.
// If no such pair exists, it returns the zero values and false. The entire
// sequence is consumed.
func (s Seq2[K, V]) LastFunc(f func(K, V) bool) (K, V, bool) {
	return xiter.LastFunc2(s.Iter(), f)
}

// Position is a terminal operation that returns the zero-based index of the
// first pair satisfying f. Returns (-1, false) when no pair matches.
// Consumption stops at the first match.
func (s Seq2[K, V]) Position(f func(K, V) bool) (int, bool) {
	return xiter.Position2(s.Iter(), f)
}

// Nth is a terminal operation that returns the n-th pair (zero-based index).
// Returns (zero, zero, false) when n is negative or when the sequence has
// fewer than n+1 pairs. Only the first n+1 pairs are consumed.
//
//	Of2(xiter.Enumerate(xiter.Range1(10))).Nth(3)  // (3, 3, true)
func (s Seq2[K, V]) Nth(n int) (K, V, bool) { return xiter.Nth2(s.Iter(), n) }

// CompareFunc is a terminal operation that lexicographically compares s and
// other pair by pair using f, following cmp.Compare's negative/zero/positive
// convention. Comparison stops at the first differing pair; if one sequence is
// a prefix of the other, the shorter sequence is considered smaller. Returns 0
// when the two sequences are equal.
func (s Seq2[K, V]) CompareFunc(other Seq2[K, V], f func(K, V, K, V) int) int {
	return xiter.CompareFunc2(s.Iter(), other.Iter(), f)
}

// EqualFunc is a terminal operation that reports whether s and other have the
// same length and every pair of corresponding entries satisfies f. Both
// sequences are fully consumed.
func (s Seq2[K, V]) EqualFunc(other Seq2[K, V], f func(K, V, K, V) bool) bool {
	return xiter.EqualFunc2(s.Iter(), other.Iter(), f)
}

// ContainsFunc is a terminal operation that reports whether any pair of s
// satisfies f. It stops as soon as f returns true. Returns false for an empty
// sequence. Equivalent to Any but named to mirror the Contains family.
func (s Seq2[K, V]) ContainsFunc(f func(K, V) bool) bool {
	return xiter.ContainsFunc2(s.Iter(), f)
}
