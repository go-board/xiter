// Package stream provides a fluent wrapper API for xiter sequences.
//
// The stream API intentionally exposes a commonly used subset of xiter
// operations for chain-style composition. It does not include every function in
// package xiter.
//
// For operations not exposed as fluent methods, use Seq/Seq2 to access the
// underlying iter.Seq/iter.Seq2 and call xiter functions directly.
package stream
