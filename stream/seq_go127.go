//go:build go1.27

package stream

import "github.com/go-board/xiter"

// Map returns a Seq[E2] obtained by applying f to each element of s. The
// output length matches the input length. Requires Go 1.27 method-level
// generics because the element type changes.
//
//	Of(xiter.Range1(3)).Map(strconv.Itoa)  // yields "0", "1", "2"
func (s Seq[E]) Map[E2 any](f func(E) E2) Seq[E2] { return Of(xiter.Map(s.Iter(), f)) }

// MapWhile applies f to each element and yields the results until f returns
// ok=false, at which point the sequence stops. The element that caused the
// stop is not yielded. Requires Go 1.27 method-level generics.
//
//	Of(xiter.Range1(10)).MapWhile(func(x int) (int, bool) {
//	    if x < 3 { return x * 10, true }
//	    return 0, false
//	})  // yields 0, 10, 20
func (s Seq[E]) MapWhile[E2 any](f func(E) (E2, bool)) Seq[E2] {
	return Of(xiter.MapWhile(s.Iter(), f))
}

// FilterMap applies f to each element and keeps only the results for which f
// returned ok=true. It is equivalent to Filter followed by Map, but in a
// single pass. Requires Go 1.27 method-level generics.
//
//	Of(seqOf(1, 2, 3, 4)).FilterMap(func(n int) (string, bool) {
//	    if n%2 == 0 {
//	        return strconv.Itoa(n), true
//	    }
//	    return "", false
//	})  // yields "2", "4"
func (s Seq[E]) FilterMap[E2 any](f func(E) (E2, bool)) Seq[E2] {
	return Of(xiter.FilterMap(s.Iter(), f))
}

// Split turns each element of s into a (key, value) pair by applying f,
// producing a Seq2[K, V]. It is the inverse direction of Join. Requires Go 1.27
// method-level generics because the output type changes to Seq2.
//
//	Of(seqOf(1, 2)).Split(func(e int) (int, int) { return e, e * 10 })
//	// yields (1,10), (2,20)
func (s Seq[E]) Split[K, V any](f func(E) (K, V)) Seq2[K, V] {
	return Of2(xiter.Split(s.Iter(), f))
}

// Zip pairs elements of s with elements of other into a Seq2[E, E2].
// Iteration stops as soon as either sequence is exhausted. Requires Go 1.27
// method-level generics because the output type changes to Seq2.
//
//	Of(seqOf(1, 2, 3)).Zip(Of(seqOf("a", "b", "c")))
//	// yields (1,"a"), (2,"b"), (3,"c")
func (s Seq[E]) Zip[E2 any](other Seq[E2]) Seq2[E, E2] {
	return Of2(xiter.Zip(s.Iter(), other.Iter()))
}

// ZipWith combines elements of s and other pairwise via f, producing a Seq[E3].
// Iteration stops as soon as either sequence is exhausted. Requires Go 1.27
// method-level generics because the output element type changes.
//
//	Of(seqOf(1, 2, 3)).ZipWith(Of(seqOf(10, 20, 30)), func(a, b int) int {
//	    return a + b
//	})  // yields 11, 22, 33
func (s Seq[E]) ZipWith[E2, E3 any](other Seq[E2], f func(E, E2) E3) Seq[E3] {
	return Of(xiter.ZipWith(s.Iter(), other.Iter(), f))
}

// Fold is a terminal operation that reduces s to a single value by repeatedly
// applying f(acc, element), starting from init. Returns init unchanged when
// the sequence is empty. Requires Go 1.27 method-level generics because the
// accumulator type is independent of E.
//
//	Of(seqOf(1, 2, 3, 4)).Fold(0, func(acc, n int) int { return acc + n })  // 10
func (s Seq[E]) Fold[A any](init A, f func(A, E) A) A { return xiter.Fold(s.Iter(), init, f) }

// TryFold is like Fold but stops when f returns an error. Returns the
// accumulator built so far together with the first error; on success returns
// (final accumulator, nil). Requires Go 1.27 method-level generics.
func (s Seq[E]) TryFold[A any](init A, f func(A, E) (A, error)) (A, error) {
	return xiter.TryFold(s.Iter(), init, f)
}

// Scan returns a Seq[A] that maintains an accumulator seeded with init and
// yields the accumulator after each element is folded in by f. The first
// yielded value reflects the accumulator after consuming the first element.
// Stops early when f returns ok=false. Requires Go 1.27 method-level generics.
//
//	Of(seqOf(1, 2, 3, 4)).Scan(0, func(acc, n int) (int, bool) {
//	    return acc + n, true
//	})  // yields 1, 3, 6, 10
func (s Seq[E]) Scan[A any](init A, f func(A, E) (A, bool)) Seq[A] {
	return Of(xiter.Scan(s.Iter(), init, f))
}
