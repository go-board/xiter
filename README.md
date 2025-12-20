# xiter

A small set of iterator helpers built on Go's `iter` package. It provides
sequence (`iter.Seq`) and key/value sequence (`iter.Seq2`) utilities inspired by
iterator patterns from other languages.

## Features

- Lazy sequence transformations (map/filter/take/skip, etc.)
- Reducers and predicates (fold, any/all, contains, size helpers)
- Conversions to slices, maps, and sets
- Support for both `Seq` and `Seq2` variants

## Installation

```sh
go get github.com/go-board/xiter
```

## Usage

```go
package main

import (
	"fmt"
	"xiter"
)

func main() {
	seq := xiter.Range1(5)
	result := xiter.ToSlice(xiter.Map(seq, func(v int) int { return v * 2 }))
	fmt.Println(result)
}
```

## Examples

### Filtering and taking

```go
seq := xiter.Range1(10)
evens := xiter.Filter(seq, func(v int) bool { return v%2 == 0 })
firstTwo := xiter.Take(evens, 2)
fmt.Println(xiter.ToSlice(firstTwo))
```

### Working with key/value sequences

```go
pairs := xiter.Enumerate(xiter.Range1(3))
mapped := xiter.Map2(pairs, func(k int, v int) (int, int) {
	return k, v * 10
})
fmt.Println(xiter.ToSlice2(mapped))
```

### Grouping values

```go
words := xiter.FromSlice([]string{"apple", "banana", "apricot", "cherry"})
groups := xiter.GroupBy(words, func(s string) byte { return s[0] })
fmt.Println(groups["a"])
```

## API Overview

Some commonly used helpers include:

- `Map`, `Filter`, `Take`, `Skip`, `TakeWhile`, `SkipWhile`
- `Fold`, `Sum`, `Any`, `All`, `Contains`
- `Find`, `Position`, `Enumerate`
- `ToSlice`, `ToMap`, `ToSet`

Key/value variants use the `*2` suffix (e.g. `Map2`, `Filter2`, `Take2`).

## License

MIT
