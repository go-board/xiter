# xiter - Extended Iterators for Go

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-board/xiter.svg)](https://pkg.go.dev/github.com/go-board/xiter)
[![Go Version](https://img.shields.io/github/go-mod/go-version/go-board/xiter)](https://github.com/go-board/xiter/blob/main/go.mod)
[![Build Status](https://github.com/go-board/xiter/actions/workflows/go.yml/badge.svg)](https://github.com/go-board/xiter/actions/workflows/go.yml)

xiter is a powerful extension library for Go's standard `iter` package, providing a rich set of functional-style operations for working with sequences and key-value pairs in Go 1.22+.

## Features

- **Comprehensive API**: Over 40 functional operations for sequence manipulation
- **Type-safe**: Full support for Go's generics throughout the API
- **Dual sequence support**: Works with both single-element sequences (`iter.Seq`) and key-value pairs (`iter.Seq2`)
- **Lazy evaluation**: All operations are lazily evaluated for optimal performance
- **Side-effect free**: Most operations return new iterators without modifying the original data
- **Seamless integration**: Drop-in compatibility with Go's standard `iter` package
- **Performance optimized**: Minimal overhead and efficient implementations

## Installation

```bash
go get github.com/go-board/xiter
```

## Quick Start

```go
package main

import (
	"fmt"
	"iter"
	"github.com/go-board/xiter"
)

func main() {
	// Generate a sequence of numbers from 0 to 9
	numbers := xiter.Range1(10)

	// Double each number
	doubled := xiter.Map(numbers, func(x int) int { return x * 2 })

	// Filter even numbers (which are all even after doubling)
	evenDoubles := xiter.Filter(doubled, func(x int) bool { return x%2 == 0 })

	// Sum the results
	sum := xiter.Fold(evenDoubles, 0, func(acc, x int) int { return acc + x })

	fmt.Printf("Sum of doubled numbers: %d\n", sum) // Output: 90

	// Working with key-value pairs
	pairs := xiter.FromFunc2(func() (string, int, bool) {
		// This would typically come from some data source
		staticPairs := []struct{ K string; V int }{"one": 1, "two": 2, "three": 3}
		index := 0
		return func() (string, int, bool) {
			if index >= len(staticPairs) {
				return "", 0, false
			}
			pair := staticPairs[index]
			index++
			return pair.K, pair.V, true
		}
	}())

	// Multiply values by 10
	enhanced := xiter.Map2(pairs, func(k string, v int) (string, int) {
		return k, v * 10
	})

	// Print all pairs
	xiter.ForEach2(enhanced, func(k string, v int) {
		fmt.Printf("%s: %d\n", k, v) // Outputs: one: 10, two: 20, three: 30
	})
}
```

## Core Concepts

xiter provides operations for two main sequence types:

- **`iter.Seq[E]`**: A sequence of elements of type `E`
- **`iter.Seq2[K, V]`**: A sequence of key-value pairs of types `K` and `V`

All operations are designed to be chained together, creating a pipeline of transformations that are only executed when the sequence is consumed.

## API Overview

### Sequence Generation

```go
// Create a sequence from 0 to 9
numbers := xiter.Range1(10)

// Create a sequence from 5 to 14
numbers := xiter.Range2(5, 15)

// Create a sequence from 1 to 10 with step 2
odds := xiter.Range3(1, 11, 2)

// Create a sequence from a function
count := 0
seq := xiter.FromFunc(func() (int, bool) {
    count++
    return count, count <= 5
})

// Create a sequence with a single element
one := xiter.Once(42)

// Create an empty sequence
empty := xiter.Empty[int]()

// Create an infinite sequence repeating a value
repeats := xiter.Repeat("hello")
```

### Mapping

```go
// Transform each element
numbers := xiter.Range1(5)
squares := xiter.Map(numbers, func(x int) int { return x * x })

// Transform while condition is true
doubles := xiter.MapWhile(numbers, func(x int) (int, bool) {
    result := x * 2
    return result, result < 10
})

// Transform key-value pairs
pairs := xiter.FromFunc2(func() (string, int, bool) {
    // ...
})
transformed := xiter.Map2(pairs, func(k string, v int) (string, int) {
    return "key_" + k, v * 10
})
```

### Filtering

```go
// Filter elements that satisfy a condition
numbers := xiter.Range1(10)
evens := xiter.Filter(numbers, func(x int) bool { return x%2 == 0 })

// Filter and map in a single step
strings := []string{"one", "two", "three"}
seq := xiter.FromFunc(func() (string, bool) {
    // ... (yield strings from slice)
})
longStrings := xiter.FilterMap(seq, func(s string) (string, bool) {
    if len(s) > 3 {
        return s, true
    }
    return "", false
})
```

### Reduction

```go
// Apply a function to each element (for side effects)
numbers := xiter.Range1(5)
xiter.ForEach(numbers, func(x int) {
    fmt.Println(x)
})

// Reduce a sequence to a single value
numbers := xiter.Range1(10)
sum := xiter.Fold(numbers, 0, func(acc, x int) int { return acc + x })
```

### Searching and Checking

```go
// Check if a sequence contains an element
numbers := xiter.Range1(10)
hasFive := xiter.Contains(numbers, 5)

// Check if any element satisfies a condition
hasEven := xiter.Any(numbers, func(x int) bool { return x%2 == 0 })

// Check if all elements satisfy a condition
allPositive := xiter.All(numbers, func(x int) bool { return x > 0 })

// Find the first element that satisfies a condition
firstEven, found := xiter.FirstFunc(numbers, func(x int) bool { return x%2 == 0 })

// Find the last element
last, found := xiter.Last(numbers)
```

### Sequence Operations

```go
// Concatenate two sequences
seq1 := xiter.Range1(3)
seq2 := xiter.Range2(5, 8)
combined := xiter.Chain(seq1, seq2)

// Take the first 5 elements
firstFive := xiter.Take(numbers, 5)

// Skip the first 3 elements
rest := xiter.Skip(numbers, 3)

// Take elements while condition is true
untilFive := xiter.TakeWhile(numbers, func(x int) bool { return x < 5 })
```

## Key-Value Pair Operations

xiter provides equivalent operations for key-value pair sequences (`iter.Seq2[K, V]`):

```go
// Create key-value pairs
pairs := xiter.FromFunc2(func() (string, int, bool) {
    // ...
})

// Transform pairs
enhanced := xiter.Map2(pairs, func(k string, v int) (string, int) {
    return "enhanced_" + k, v * 2
})

// Filter pairs
validPairs := xiter.Filter2(enhanced, func(k string, v int) bool {
    return v > 5
})

// Reduce pairs
sum := xiter.Fold2(validPairs, 0, func(acc int, k string, v int) int {
    return acc + v
})
```

## Best Practices

1. **Chain operations for readability**: Combine multiple operations in a single pipeline
2. **Consume sequences efficiently**: Iterators can only be consumed once
3. **Use lazy evaluation**: Take advantage of xiter's lazy evaluation for large datasets
4. **Avoid unnecessary materialization**: Prefer to work with iterators directly rather than converting to slices
5. **Leverage type inference**: Let Go's type inference work for you when possible

## Performance Considerations

- All xiter operations are designed to be lightweight and efficient
- Lazy evaluation means only elements that are actually consumed are processed
- Memory usage is minimized since sequences are not fully materialized unless explicitly requested
- For maximum performance, avoid unnecessary conversions between sequences and concrete types

## Examples

### Calculating Factorial

```go
func Factorial(n int) int {
    numbers := xiter.Range2(1, n+1)
    return xiter.Fold(numbers, 1, func(acc, x int) int { return acc * x })
}
```

### Processing Data from a Database

```go
func ProcessUsers(db *sql.DB) error {
    rows, err := db.Query("SELECT id, name, email FROM users")
    if err != nil {
        return err
    }
    defer rows.Close()

    // Create a sequence from database rows
    userSeq := xiter.FromFunc2(func() (int, string, bool) {
        if !rows.Next() {
            return 0, "", false
        }
        var id int
        var name string
        if err := rows.Scan(&id, &name, nil); err != nil {
            return 0, "", false
        }
        return id, name, true
    })

    // Filter active users (assuming some condition)
    activeUsers := xiter.Filter2(userSeq, func(id int, name string) bool {
        // Implement your active user logic
        return id%2 == 0 // Example condition
    })

    // Process each active user
    xiter.ForEach2(activeUsers, func(id int, name string) {
        fmt.Printf("Processing user: %s (ID: %d)\n", name, id)
    })

    return nil
}
```

### Working with Large Datasets

```go
func ProcessLargeFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Create a sequence from file lines
    lines := xiter.FromFunc(func() (string, bool) {
        if !scanner.Scan() {
            return "", false
        }
        return scanner.Text(), true
    })

    // Filter and process lines
    longLines := xiter.Filter(lines, func(line string) bool { return len(line) > 100 })
    count := xiter.Fold(longLines, 0, func(acc int, line string) int {
        return acc + 1
    })

    fmt.Printf("Found %d lines longer than 100 characters\n", count)
    return nil
}
```

## Contributing

Contributions are welcome! Please feel free to submit issues, pull requests, or feature suggestions.

### Development Setup

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Run tests: `go test ./...`
6. Submit a pull request

## License

xiter is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by functional programming languages like Haskell, Scala, and JavaScript
- Built on Go's excellent new iterators feature
- Thanks to the Go community for their feedback and contributions

## Support

If you encounter any issues or have questions, please file an issue on the [GitHub repository](https://github.com/go-board/xiter/issues).

## Related Projects

- [Go Standard `iter` Package](https://pkg.go.dev/iter)
- [Go Generics](https://go.dev/doc/tutorial/generics)
