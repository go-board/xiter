package stream

import (
	"iter"

	"github.com/go-board/xiter"
)

// Stream wraps iter.Seq and provides chainable stream-style methods.
//
// Stream exposes a curated subset of xiter operations. Use Seq to access the
// underlying sequence when you need xiter functions that are not wrapped as
// methods.
type Stream[E any] struct{ seq iter.Seq[E] }

// Stream2 wraps iter.Seq2 and provides chainable stream-style methods.
//
// Stream2 exposes a curated subset of xiter key/value operations. Use Seq2 to
// access the underlying sequence when you need xiter functions that are not
// wrapped as methods.
type Stream2[K, V any] struct{ seq iter.Seq2[K, V] }

// Of wraps an iter.Seq as a Stream.
func Of[E any](seq iter.Seq[E]) Stream[E] { return Stream[E]{seq: seq} }

// Of2 wraps an iter.Seq2 as a Stream2.
func Of2[K, V any](seq iter.Seq2[K, V]) Stream2[K, V] { return Stream2[K, V]{seq: seq} }

// Seq returns the underlying iter.Seq.
func (s Stream[E]) Seq() iter.Seq[E] { return s.seq }

// Seq2 returns the underlying iter.Seq2.
func (s Stream2[K, V]) Seq2() iter.Seq2[K, V] { return s.seq }

// Filter keeps elements that match the predicate.
func (s Stream[E]) Filter(f func(E) bool) Stream[E] { return Of(xiter.Filter(s.seq, f)) }

// Map maps each element to another element of the same type.
func (s Stream[E]) Map(f func(E) E) Stream[E] { return Of(xiter.Map(s.seq, f)) }

// MapWhile maps elements until the callback returns ok=false.
func (s Stream[E]) MapWhile(f func(E) (E, bool)) Stream[E] { return Of(xiter.MapWhile(s.seq, f)) }

// FilterMap maps elements and keeps only mapped results with ok=true.
func (s Stream[E]) FilterMap(f func(E) (E, bool)) Stream[E] { return Of(xiter.FilterMap(s.seq, f)) }

// FlatMap maps each element to a sub-sequence and flattens them.
func (s Stream[E]) FlatMap(f func(E) iter.Seq[E]) Stream[E] { return Of(xiter.FlatMap(s.seq, f)) }

// Take returns a stream of the first n elements.
func (s Stream[E]) Take(n int) Stream[E] { return Of(xiter.Take(s.seq, n)) }

// Skip returns a stream after skipping the first n elements.
func (s Stream[E]) Skip(n int) Stream[E] { return Of(xiter.Skip(s.seq, n)) }

// TakeWhile returns elements while predicate is true.
func (s Stream[E]) TakeWhile(f func(E) bool) Stream[E] { return Of(xiter.TakeWhile(s.seq, f)) }

// SkipWhile skips elements while predicate is true.
func (s Stream[E]) SkipWhile(f func(E) bool) Stream[E] { return Of(xiter.SkipWhile(s.seq, f)) }

// Chain concatenates the current stream with another stream.
func (s Stream[E]) Chain(other Stream[E]) Stream[E] { return Of(xiter.Chain(s.seq, other.seq)) }

// Enumerate converts Stream[E] into Stream2[index, value].
func (s Stream[E]) Enumerate() Stream2[int, E] { return Of2(xiter.Enumerate(s.seq)) }

// ForEach consumes the stream and calls f for each element.
func (s Stream[E]) ForEach(f func(E)) { xiter.ForEach(s.seq, f) }

// Fold reduces elements into one value of the same type.
func (s Stream[E]) Fold(init E, f func(E, E) E) E { return xiter.Fold(s.seq, init, f) }

// Size returns the number of elements in the stream.
func (s Stream[E]) Size() int { return xiter.Size(s.seq) }

// SizeFunc counts elements matching predicate f.
func (s Stream[E]) SizeFunc(f func(E) bool) int { return xiter.SizeFunc(s.seq, f) }

// Any reports whether any element satisfies the predicate.
func (s Stream[E]) Any(f func(E) bool) bool { return xiter.Any(s.seq, f) }

// All reports whether all elements satisfy the predicate.
func (s Stream[E]) All(f func(E) bool) bool { return xiter.All(s.seq, f) }

// First returns the first element and whether it exists.
func (s Stream[E]) First() (E, bool) { return xiter.First(s.seq) }

// Last returns the last element and whether it exists.
func (s Stream[E]) Last() (E, bool) { return xiter.Last(s.seq) }

// FirstFunc returns the first element matching predicate f.
func (s Stream[E]) FirstFunc(f func(E) bool) (E, bool) { return xiter.FirstFunc(s.seq, f) }

// LastFunc returns the last element matching predicate f.
func (s Stream[E]) LastFunc(f func(E) bool) (E, bool) { return xiter.LastFunc(s.seq, f) }

// Position returns the index of the first matching element.
func (s Stream[E]) Position(f func(E) bool) (int, bool) { return xiter.Position(s.seq, f) }

// IsSortedFunc reports whether stream is sorted by comparator f.
func (s Stream[E]) IsSortedFunc(f func(E, E) int) bool { return xiter.IsSortedFunc(s.seq, f) }

// CompareFunc compares two streams with comparator f.
func (s Stream[E]) CompareFunc(other Stream[E], f func(E, E) int) int {
	return xiter.CompareFunc(s.seq, other.seq, f)
}

// EqualFunc reports whether two streams are equal by predicate f.
func (s Stream[E]) EqualFunc(other Stream[E], f func(E, E) bool) bool {
	return xiter.EqualFunc(s.seq, other.seq, f)
}

// MaxFunc returns max element by comparator f.
func (s Stream[E]) MaxFunc(f func(E, E) int) (E, bool) { return xiter.MaxFunc(s.seq, f) }

// MinFunc returns min element by comparator f.
func (s Stream[E]) MinFunc(f func(E, E) int) (E, bool) { return xiter.MinFunc(s.seq, f) }

// ContainsFunc reports whether any element satisfies predicate f.
func (s Stream[E]) ContainsFunc(f func(E) bool) bool { return xiter.ContainsFunc(s.seq, f) }

// Filter keeps key/value pairs that match the predicate.
func (s Stream2[K, V]) Filter(f func(K, V) bool) Stream2[K, V] { return Of2(xiter.Filter2(s.seq, f)) }

// Map maps each key/value pair to another pair of the same key/value types.
func (s Stream2[K, V]) Map(f func(K, V) (K, V)) Stream2[K, V] { return Of2(xiter.Map2(s.seq, f)) }

// FilterMap maps pairs and keeps only results with ok=true.
func (s Stream2[K, V]) FilterMap(f func(K, V) (K, V, bool)) Stream2[K, V] {
	return Of2(xiter.FilterMap2(s.seq, f))
}

// MapWhile maps pairs until callback returns ok=false.
func (s Stream2[K, V]) MapWhile(f func(K, V) (K, V, bool)) Stream2[K, V] {
	return Of2(xiter.MapWhile2(s.seq, f))
}

// Take returns a stream with first n key/value pairs.
func (s Stream2[K, V]) Take(n int) Stream2[K, V] { return Of2(xiter.Take2(s.seq, n)) }

// Skip returns a stream after skipping first n key/value pairs.
func (s Stream2[K, V]) Skip(n int) Stream2[K, V] { return Of2(xiter.Skip2(s.seq, n)) }

// TakeWhile returns pairs while predicate is true.
func (s Stream2[K, V]) TakeWhile(f func(K, V) bool) Stream2[K, V] {
	return Of2(xiter.TakeWhile2(s.seq, f))
}

// SkipWhile skips pairs while predicate is true.
func (s Stream2[K, V]) SkipWhile(f func(K, V) bool) Stream2[K, V] {
	return Of2(xiter.SkipWhile2(s.seq, f))
}

// Chain concatenates the current Stream2 with another Stream2.
func (s Stream2[K, V]) Chain(other Stream2[K, V]) Stream2[K, V] {
	return Of2(xiter.Chain2(s.seq, other.seq))
}

// ForEach consumes the stream and calls f for each key/value pair.
func (s Stream2[K, V]) ForEach(f func(K, V)) { xiter.ForEach2(s.seq, f) }

// Size returns the number of key/value pairs.
func (s Stream2[K, V]) Size() int { return xiter.Size2(s.seq) }

// SizeFunc counts key/value pairs matching predicate f.
func (s Stream2[K, V]) SizeFunc(f func(K, V) bool) int { return xiter.SizeFunc2(s.seq, f) }

// Any reports whether any pair satisfies predicate f.
func (s Stream2[K, V]) Any(f func(K, V) bool) bool { return xiter.Any2(s.seq, f) }

// All reports whether all pairs satisfy predicate f.
func (s Stream2[K, V]) All(f func(K, V) bool) bool { return xiter.All2(s.seq, f) }

// First returns the first pair and whether it exists.
func (s Stream2[K, V]) First() (K, V, bool) { return xiter.First2(s.seq) }

// Last returns the last pair and whether it exists.
func (s Stream2[K, V]) Last() (K, V, bool) { return xiter.Last2(s.seq) }

// FirstFunc returns the first key/value pair matching predicate f.
func (s Stream2[K, V]) FirstFunc(f func(K, V) bool) (K, V, bool) { return xiter.FirstFunc2(s.seq, f) }

// LastFunc returns the last key/value pair matching predicate f.
func (s Stream2[K, V]) LastFunc(f func(K, V) bool) (K, V, bool) { return xiter.LastFunc2(s.seq, f) }

// Position returns the index of first pair matching predicate f.
func (s Stream2[K, V]) Position(f func(K, V) bool) (int, bool) { return xiter.Position2(s.seq, f) }

// CompareFunc compares two Stream2 by comparator f.
func (s Stream2[K, V]) CompareFunc(other Stream2[K, V], f func(K, V, K, V) int) int {
	return xiter.CompareFunc2(s.seq, other.seq, f)
}

// EqualFunc reports whether two Stream2 are equal by predicate f.
func (s Stream2[K, V]) EqualFunc(other Stream2[K, V], f func(K, V, K, V) bool) bool {
	return xiter.EqualFunc2(s.seq, other.seq, f)
}

// ContainsFunc reports whether any key/value pair satisfies predicate f.
func (s Stream2[K, V]) ContainsFunc(f func(K, V) bool) bool { return xiter.ContainsFunc2(s.seq, f) }
