# xiter

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-board/xiter.svg)](https://pkg.go.dev/github.com/go-board/xiter)
[![Go Version](https://img.shields.io/github/go-mod/go-version/go-board/xiter)](https://github.com/go-board/xiter/blob/main/go.mod)
[![Build Status](https://github.com/go-board/xiter/actions/workflows/go.yml/badge.svg)](https://github.com/go-board/xiter/actions/workflows/go.yml)

`xiter` is an extension toolkit for Go's standard `iter` package. It provides functional, lazy, and type-safe sequence operations for both `iter.Seq[E]` and `iter.Seq2[K, V]`.

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Functional API (`xiter`)](#functional-api-xiter)
  - [Fluent API (`stream` subpackage)](#fluent-api-stream-subpackage)
- [Core Concepts](#core-concepts)
- [API Overview](#api-overview)
  - [Source](#source)
  - [Transform](#transform)
  - [Filter / Slice](#filter--slice)
  - [Terminal](#terminal)
  - [Compare / Search](#compare--search)
  - [`stream` subpackage](#stream-subpackage)
  - [`collector` subpackage (experimental)](#collector-subpackage-experimental)
- [Best Practices](#best-practices)
- [Development & Testing](#development--testing)
- [Roadmap](#roadmap)
- [License](#license)

---

## Features

- ✅ Supports both `iter.Seq[E]` and `iter.Seq2[K, V]`
- ✅ Fully generic and type-safe
- ✅ Lazy evaluation with on-demand consumption
- ✅ Covers common workflows: map, filter, slice, reduce, search, compare, sorted checks
- ✅ Supports error-aware terminal operations (`TryForEach`, `TryFold`, `TryReduce`)
- ✅ Provides both functional APIs (`xiter`) and fluent APIs (`xiter/stream`)
- ✅ Includes source generators: `Range`, `Iterate`, `FromFunc`, `Once`, `Empty`, `Repeat`
- 🧪 Experimental `collector` subpackage with reusable terminal collectors (`ToSlice`, `ToMap`, `GroupingBy`, ...)

## Installation

```bash
go get github.com/go-board/xiter
```

## Quick Start

### Functional API (`xiter`)

```go
package main

import (
	"fmt"

	"github.com/go-board/xiter"
)

func main() {
	numbers := xiter.Range1(10) // 0..9
	evens := xiter.Filter(numbers, func(v int) bool { return v%2 == 0 })
	doubled := xiter.Map(evens, func(v int) int { return v * 2 })

	sum := xiter.Fold(doubled, 0, func(acc, v int) int { return acc + v })
	fmt.Println(sum) // 40
}
```

### Fluent API (`stream` subpackage)

> `stream` provides chainable `Seq` / `Seq2` function types. Core same-type
> methods work without method-level generics. Methods that need their own type
> parameters, such as `Map`, `Fold`, `Split`, `Join`, `Zip`, and `ZipWith`,
> require Go 1.27 or newer.

```go
package main

import (
	"fmt"

	"github.com/go-board/xiter"
	"github.com/go-board/xiter/stream"
)

func main() {
	s := stream.Of(xiter.Range1(10)).
		Skip(2).
		Take(5).
		Filter(func(v int) bool { return v%2 == 0 })

	result := make([]int, 0)
	s.ForEach(func(v int) { result = append(result, v*10) })

	fmt.Println(result) // [20 40 60]
}
```

With Go 1.27 generic methods, `stream` also supports type-changing chains:

```go
s := stream.Of(xiter.Range1(5)).
	Map(func(v int) string { return fmt.Sprintf("n=%d", v) }).
	Fold([]string{}, func(acc []string, v string) []string {
		return append(acc, v)
	})

fmt.Println(s) // [n=0 n=1 n=2 n=3 n=4]
```

`Iterate` generates a sequence from a seed and a step function that returns
`(next, ok)`; the sequence stops when `ok` is false, so no external limiter is
needed:

```go
// 1, 2, 4, 8, 16 — stops once x reaches 16.
stream.Iterate(1, func(x int) (int, bool) {
	if x >= 16 {
		return 0, false
	}
	return x * 2, true
}).ForEach(func(v int) { fmt.Println(v) })
```

## Core Concepts

- `iter.Seq[E]`: sequence of single values
- `iter.Seq2[K, V]`: sequence of key/value pairs
- `stream.Seq[E]` and `stream.Seq2[K, V]`: chainable wrappers over the bare iterator function types (`func(yield func(E) bool)` / `func(yield func(K, V) bool)`), the same underlying type as `iter.Seq` / `iter.Seq2`
- Sequences are **lazy**: execution happens at terminal stages like `ForEach`, `Fold`, `First`, and `Last`
- Sequences are usually **single-pass**: avoid re-consuming the same exhausted source

## API Overview

> See full signatures on GoDoc: <https://pkg.go.dev/github.com/go-board/xiter>

### Source

- `Range1`, `Range2`, `Range3`
- `FromFunc`, `FromFunc2`
- `Iterate`, `Iterate2`
- `Once`, `Once2`
- `Empty`, `Empty2`
- `Repeat`, `Repeat2`

### Transform

- `Map`, `Map2`
- `MapWhile`, `MapWhile2`
- `FlatMap`, `Flatten`
- `Inspect`, `Inspect2`
- `Enumerate`
- `Join`, `Split`
- `Keys`, `Values`, `Swap`
- `Cast`
- `Scan`

### Filter / Slice

- `Filter`, `Filter2`
- `FilterMap`, `FilterMap2`
- `Take`, `Take2`, `TakeWhile`, `TakeWhile2`
- `Skip`, `Skip2`, `SkipWhile`, `SkipWhile2`
- `Chain`, `Chain2`
- `Zip`, `ZipWith`

### Terminal

- `ForEach`, `ForEach2`
- `TryForEach`, `TryForEach2`
- `Fold`, `Fold2`, `TryFold`, `TryFold2`
- `Reduce`, `Reduce2`, `TryReduce`, `TryReduce2`
- `Size`, `Size2`, `SizeFunc`, `SizeFunc2`, `SizeValue`, `SizeValue2`

### Compare / Search

- `Contains`, `Contains2`, `ContainsFunc`, `ContainsFunc2`
- `Any`, `Any2`, `All`, `All2`
- `First`, `First2`, `FirstFunc`, `FirstFunc2`
- `Last`, `Last2`, `LastFunc`, `LastFunc2`
- `Position`, `Position2`
- `Compare`, `Compare2`, `CompareFunc`, `CompareFunc2`
- `Equal`, `Equal2`, `EqualFunc`, `EqualFunc2`
- `Max`, `MaxFunc`, `Min`, `MinFunc`
- `MinMax`, `MinMaxFunc`
- `IsSorted`, `IsSortedFunc`

### `stream` subpackage

The `stream` subpackage exposes chainable function types:

- `stream.Seq[E]`
- `stream.Seq2[K, V]`
- `stream.Of`, `stream.Of2` — wrap a bare iterator function (`func(yield func(E) bool)` / `func(yield func(K, V) bool)`, the same underlying type as `iter.Seq` / `iter.Seq2`)
- `stream.FromFunc`, `stream.FromFunc2`
- `stream.Iterate`, `stream.Iterate2`
- `Iter` to get back the underlying `iter.Seq` or `iter.Seq2`

Available without Go 1.27 method-level generics:

- `Seq`: `Filter`, `Inspect`, `Take`, `Skip`, `TakeWhile`, `SkipWhile`, `Chain`, `Enumerate`
- `Seq`: `ForEach`, `TryForEach`, `Reduce`, `TryReduce`
- `Seq`: `Size`, `SizeFunc`, `Any`, `All`, `First`, `Last`, `FirstFunc`, `LastFunc`, `Position`
- `Seq`: `IsSortedFunc`, `CompareFunc`, `EqualFunc`, `MaxFunc`, `MinFunc`, `MinMaxFunc`, `ContainsFunc`
- `Seq2`: `Filter`, `Keys`, `Values`, `Swap`, `Inspect`, `Take`, `Skip`, `TakeWhile`, `SkipWhile`, `Chain`
- `Seq2`: `ForEach`, `TryForEach`, `Reduce`, `TryReduce`
- `Seq2`: `Size`, `SizeFunc`, `Any`, `All`, `First`, `Last`, `FirstFunc`, `LastFunc`, `Position`
- `Seq2`: `CompareFunc`, `EqualFunc`, `ContainsFunc`

Available when building with Go 1.27 or newer:

- `Seq`: `Map`, `MapWhile`, `FilterMap`, `Split`, `Zip`, `ZipWith`, `Fold`, `TryFold`, `Scan`
- `Seq2`: `Map`, `MapWhile`, `FilterMap`, `Join`, `Fold`, `TryFold`

### `collector` subpackage (experimental)

> ⚠️ Experimental: the `collector` API is not yet stable and may change
> incompatibly or be removed in a future version.

The `collector` subpackage provides reusable terminal operations that
materialize an `iter.Seq` / `iter.Seq2` into a container or aggregated value.
A `Collector[E, R]` is a named, reusable function; collectors compose, e.g.
`GroupingByDownstream` delegates to a downstream collector per group.

```go
s := xiter.Range1(10)
got := collector.Collect(s, collector.ToSlice[int]())
// got == []int{0, 1, 2, ..., 9}

// group even/odd and materialize each group
groups := collector.Collect(
    xiter.Range1(10),
    collector.GroupingByDownstream(
        func(n int) int { return n % 2 },
        collector.ToSlice[int](),
    ),
)
// groups == map[int][]int{0: {0,2,4,6,8}, 1: {1,3,5,7,9}}
```

Core types:

- `Collector[E, R]`, `Collector2[K, V, R]`
- `Collect`, `Collect2`

Collectors for `iter.Seq[E]`:

- `ToSlice`, `ToSet`
- `ToMap`, `ToMapMerge`
- `Joining`
- `GroupingBy`, `GroupingByDownstream`
- `PartitioningBy` (returns `Partition[E]{Pass, Fail}`)

Collectors for `iter.Seq2[K, V]`:

- `ToMap2`, `ToMap2Merge`
- `ToKeys`, `ToValues`

## Best Practices

1. Compose transformations as pipelines for readability.
2. Delay materialization (e.g. slice/map conversion) whenever possible.
3. `Repeat` produces an infinite sequence — always pair it with a limiting operator such as `Take`. `Iterate` is self-terminating: its `next` callback returns `(value, ok)` and stops when `ok` is false.
4. Use `*_Func` variants for custom comparison and matching logic.
5. Use `Inspect` for debugging or side effects in the middle of a lazy pipeline.
6. Use `Try*` terminal operations when callbacks can fail and should stop early.
7. Honor the `yield` return value in custom sources: stop calling `yield` as soon as it returns `false`, otherwise Go 1.23+ range-over-func will panic.

## Development & Testing

```bash
go test ./...
```

## Roadmap

- Add more examples and benchmarks
- Continue improving `stream` API coverage and documentation

## License

Apache-2.0. See [LICENSE](LICENSE) for details.
