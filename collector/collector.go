// Package collector provides reusable, composable terminal operations — called
// collectors — that materialize an iter.Seq or iter.Seq2 into a container or
// aggregated value.
//
// A Collector is a named function value that consumes a sequence and produces a
// result. Unlike the terminal operations in the root xiter package, which are
// one-shot consumers inlined at the call site, collectors are first-class
// values that can be stored, passed around, and composed: GroupingByDownstream,
// for example, applies a downstream collector to each group, so the same
// ToSlice or counting collector can be reused inside grouping.
//
// Apply a collector either by calling it directly or through Collect, which
// reads as the conventional "collect a sequence into a container" verb:
//
//	s := xiter.Range1(5)
//	got := collector.Collect(s, collector.ToSlice[int]())
//	// got == []int{0, 1, 2, 3, 4}
//
// EXPERIMENTAL: this package is experimental. Its API is not yet stable and may
// change incompatibly or be removed in a future version.
package collector

import "iter"

// Collector is a reusable terminal operation over iter.Seq[E]: it is a function
// that consumes a sequence and produces a value of type R. It is an alias of
// func(iter.Seq[E]) R, so it may be applied either by direct call (c(s)) or via
// Collect. Collectors compose — GroupingByDownstream accepts a downstream
// Collector to materialize each group.
//
// EXPERIMENTAL: the Collector API may change in future versions.
type Collector[E, R any] func(iter.Seq[E]) R

// Collector2 is the iter.Seq2[K, V] variant of Collector: a function that
// consumes a key/value sequence and produces a value of type R. It is an alias
// of func(iter.Seq2[K, V]) R and may be applied by direct call or via Collect2.
//
// EXPERIMENTAL: the Collector2 API may change in future versions.
type Collector2[K, V, R any] func(iter.Seq2[K, V]) R

// Collect applies c to s and returns the result. It is equivalent to c(s) but
// provides a symmetric, readable API matching Collect2.
func Collect[E, R any](s iter.Seq[E], c Collector[E, R]) R { return c(s) }

// Collect2 applies c to s and returns the result. It is the iter.Seq2 variant
// of Collect and is equivalent to c(s).
func Collect2[K, V, R any](s iter.Seq2[K, V], c Collector2[K, V, R]) R { return c(s) }
