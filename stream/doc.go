// Package stream provides fluent, chainable wrappers around the functions in
// package xiter for the bare iterator function types func(yield func(E) bool)
// and func(yield func(K, V) bool) (the same underlying types as iter.Seq and
// iter.Seq2).
//
// The Seq[E] and Seq2[K, V] types wrap the underlying iterator function types
// and expose the same lazy, single-pass operations as methods, so that
// pipelines can be written in a fluent style instead of nested function calls.
// Every non-terminal method returns a new Seq/Seq2 that defers all work until
// iteration; terminal methods consume the sequence and produce a value. All
// operations honor the yield-return-false early stop signal and release the
// source as soon as the consumer breaks early.
//
// Construct a stream with Of/Of2, FromFunc/FromFunc2, or Iterate/Iterate2;
// chain transform and filter methods; terminate with a method such as ForEach,
// Reduce, First, or Fold.
//
// Methods whose element type changes — Map, MapWhile, FilterMap, Split, Zip,
// ZipWith, Join, Fold, TryFold, Scan — require Go 1.27 method-level generics
// and are only available when building with Go 1.27 or newer.
package stream
