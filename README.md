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
- ✅ Provides both functional APIs (`xiter`) and fluent APIs (`xiter/stream`)

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

> `stream` provides a fluent wrapper (`Stream` / `Stream2`) for common workflows.
> It intentionally exposes a subset of `xiter` APIs rather than a 1:1 surface.

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
		Filter(func(v int) bool { return v%2 == 0 }).
		Map(func(v int) int { return v * 10 })

	result := make([]int, 0)
	s.ForEach(func(v int) { result = append(result, v) })

	fmt.Println(result) // [20 40 60]
}
```

## Core Concepts

- `iter.Seq[E]`: sequence of single values
- `iter.Seq2[K, V]`: sequence of key/value pairs
- Sequences are **lazy**: execution happens at terminal stages like `ForEach`, `Fold`, `First`, and `Last`
- Sequences are usually **single-pass**: avoid re-consuming the same exhausted source

## API Overview

> See full signatures on GoDoc: <https://pkg.go.dev/github.com/go-board/xiter>

### Source

- `Range1`, `Range2`, `Range3`
- `FromFunc`, `FromFunc2`
- `Once`, `Once2`
- `Empty`, `Empty2`
- `Repeat`, `Repeat2`

### Transform

- `Map`, `Map2`
- `MapWhile`, `MapWhile2`
- `FlatMap`, `Flatten`
- `Enumerate`
- `Join`, `Split`
- `Cast`

### Filter / Slice

- `Filter`, `Filter2`
- `FilterMap`, `FilterMap2`
- `Take`, `Take2`, `TakeWhile`, `TakeWhile2`
- `Skip`, `Skip2`, `SkipWhile`, `SkipWhile2`
- `Chain`, `Chain2`

### Terminal

- `ForEach`, `ForEach2`
- `Fold`, `Fold2`
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
- `IsSorted`, `IsSortedFunc`

## Best Practices

1. Compose transformations as pipelines for readability.
2. Delay materialization (e.g. slice/map conversion) whenever possible.
3. Always pair infinite sources (like `Repeat`) with limiting operators (`Take`, etc.).
4. Use `*_Func` variants for custom comparison and matching logic.

## Development & Testing

```bash
go test ./...
```

## Roadmap

- Add more examples and benchmarks
- Continue improving `stream` API coverage and documentation

## License

Apache-2.0. See [LICENSE](LICENSE) for details.
