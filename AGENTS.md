# AGENTS.md

Guidance for AI agents working on the `xiter` repository.

## Project Overview

`xiter` is an extension toolkit for Go's standard `iter` package, providing
functional-style, lazy, type-safe sequence operations for `iter.Seq[E]` and
`iter.Seq2[K, V]`. A fluent `stream` subpackage wraps the same operations as
chainable methods.

- Module: `github.com/go-board/xiter`
- Minimum Go version: 1.23 (see `go.mod`)
- CI Go version: 1.24 (see `.github/workflows/go.yml`)
- License: Apache-2.0

## Repository Layout

```
xiter/
├── seq.go              # All Seq[E] operations (Source/Transform/Filter/Terminal/Compare)
├── seq2.go             # All Seq2[K,V] operations (mirrors Seq with *2 suffix)
├── types.go            # Internal constraints (e.g. integral)
├── doc.go              # Package doc comment
├── seq_test.go         # All Seq tests
├── seq2_test.go        # All Seq2 tests
├── help_test.go        # Shared test helpers (ToSlice/ToMap/FromSlice/FromMap)
└── stream/             # Fluent API subpackage
    ├── seq.go          # Seq[E] methods (no method-level generics, works on Go 1.23+)
    ├── seq2.go         # Seq2[K,V] methods (no method-level generics)
    ├── seq_go127.go    # Seq[E] methods requiring Go 1.27 method-level generics
    ├── seq2_go127.go   # Seq2[K,V] methods requiring Go 1.27 method-level generics
    └── doc.go          # Subpackage doc comment
```

### File organization rules

1. **Split by input type, not by category.** All `iter.Seq[E]` operations live in
   `seq.go`; all `iter.Seq2[K, V]` operations live in `seq2.go`. Each file uses
   section comments (`// Source`, `// Transform`, `// Filter / Slice`, `// Terminal`,
   `// Compare / Search`) to group functions by README category.

2. **Cross-type converters go by input side.** Functions that convert between
   Seq and Seq2 are placed according to their *input* type:
   - `Enumerate`, `Split`, `Cast`, `Zip`, `ZipWith` (Seq → Seq2) → `seq.go`
   - `Join`, `Keys`, `Values` (Seq2 → Seq) → `seq2.go`

3. **Tests mirror source files.** `seq_test.go` tests Seq operations;
   `seq2_test.go` tests Seq2 operations. Shared helpers stay in `help_test.go`.

4. **stream subpackage split by Go version.** Methods needing method-level
   generics (e.g. `Map`, `Fold`, `Scan`) are isolated in `*_go127.go` files;
   same-type methods work on Go 1.23+ and live in `seq.go`/`seq2.go`.

## Build & Test

```bash
# Run all tests (requires Go 1.23+; stream subpackage generics need Go 1.27+)
go test ./...

# Run with coverage
go test -coverprofile=cover.out .
go tool cover -func=cover.out

# Vet
go vet ./...
```

The project targets **100% statement coverage** for the root `xiter` package.
When adding new functions, add tests covering:
- Normal path
- Empty sequence path
- Early termination (`!yield` branch) — use the `stopEarly`/`stopEarly2` helpers
- Error path for `Try*` variants

## Coding Conventions

### Function signatures

- Seq operations: `func Name[E any](s iter.Seq[E], ...) ...`
- Seq2 operations: `func Name2[K, V any](s iter.Seq2[K, V], ...) ...`
- Ordered variants use `cmp.Ordered`; custom-comparison variants use a `Func` suffix
  and a `func(E, E) int` comparator following `cmp.Compare` convention
- Functions returning `(value, ok bool)` should return the zero value on `ok=false`

### Lazy evaluation

All non-terminal operations return `iter.Seq`/`iter.Seq2` and must:
1. Be lazy — do not consume the source until the returned sequence is iterated
2. Honor `yield` returning `false` by stopping immediately and releasing the source
3. Be single-pass — never assume the source can be re-iterated

### Pull-based internal pattern

When a function needs to draw from two sources in lockstep (e.g. `Zip`,
`CompareFunc`), use `iter.Pull`/`iter.Pull2` with `defer stop()`:

```go
it, stop := iter.Pull(s)
defer stop()
for {
    e, ok := it()
    if !ok { return }
    // ...
}
```

### Testing source construction

**Critical:** Inline test sources must check `yield`'s return value. Go 1.23+
panics if a range-over-func iterator continues calling `yield` after it returned
`false`. Prefer the shared helpers:

```go
// Good — uses helper that honors stop signal
src := seqOf(1, 2, 3)

// Good — manual source checks return value
src := func(yield func(int) bool) {
    for _, v := range []int{1, 2, 3} {
        if !yield(v) { return }
    }
}

// Bad — panics when terminal ops break early
src := func(yield func(int) bool) {
    yield(1) // ignores return value
    yield(2)
}
```

### Documentation

- Every exported function has a Go doc comment starting with the function name
- Comments describe behavior, empty-sequence semantics, and termination conditions
- Update `README.md` API overview when adding new exported functions

## Adding a new operation

1. Add the Seq version to `seq.go` under the appropriate section comment
2. Add the Seq2 version (with `2` suffix) to `seq2.go` under the matching section
3. Add the `stream` wrapper(s):
   - Same-type method → `stream/seq.go` or `stream/seq2.go`
   - Type-changing method → `stream/seq_go127.go` or `stream/seq2_go127.go`
4. Add tests to `seq_test.go` / `seq2_test.go` covering all branches
5. Update `README.md` API overview list

## Dependencies

- Standard library only (`iter`, `cmp`). No external dependencies.
- `stream` subpackage imports `xiter` as its only dependency.
