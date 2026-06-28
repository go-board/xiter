//go:build go1.27

package stream

import "github.com/go-board/xiter"

// Map applies f to each key/value pair of s and yields the resulting
// (K2, V2) pairs. The output length matches the input length. Requires Go 1.27
// method-level generics because the key and value types change.
//
//	Of2(xiter.Enumerate(xiter.Range1(3))).
//	    Map(func(i, n int) (string, int) { return strconv.Itoa(i), n })
//	// yields ("0",0), ("1",1), ("2",2)
func (s Seq2[K, V]) Map[K2, V2 any](f func(K, V) (K2, V2)) Seq2[K2, V2] {
	return Of2(xiter.Map2(s.Iter(), f))
}

// MapWhile applies f to each pair and yields the results until f returns
// ok=false, at which point the sequence stops. The pair that caused the stop
// is not yielded. Requires Go 1.27 method-level generics.
//
//	Of2(xiter.Enumerate(xiter.Range1(5))).MapWhile(func(i, n int) (int, int, bool) {
//	    if i < 3 { return i, n * 10, true }
//	    return 0, 0, false
//	})  // yields (0,0), (1,10), (2,20)
func (s Seq2[K, V]) MapWhile[K2, V2 any](f func(K, V) (K2, V2, bool)) Seq2[K2, V2] {
	return Of2(xiter.MapWhile2(s.Iter(), f))
}

// FilterMap applies f to each pair and keeps only the results for which f
// returned ok=true. It is equivalent to Filter2 followed by Map2, but in a
// single pass. Requires Go 1.27 method-level generics.
//
//	Of2(xiter.Enumerate(xiter.Range1(4))).FilterMap(func(i, n int) (string, int, bool) {
//	    if i%2 == 0 { return strconv.Itoa(i), n, true }
//	    return "", 0, false
//	})  // yields ("0",0), ("2",2)
func (s Seq2[K, V]) FilterMap[K2, V2 any](f func(K, V) (K2, V2, bool)) Seq2[K2, V2] {
	return Of2(xiter.FilterMap2(s.Iter(), f))
}

// Join turns each key/value pair of s into a single element by applying f,
// producing a Seq[E]. It is the inverse direction of Split. Requires Go 1.27
// method-level generics because the output type changes to Seq.
//
//	Of2(xiter.Enumerate(xiter.Range1(3))).
//	    Join(func(i, n int) string { return fmt.Sprintf("%d:%d", i, n) })
//	// yields "0:0", "1:1", "2:2"
func (s Seq2[K, V]) Join[E any](f func(K, V) E) Seq[E] { return Of(xiter.Join(s.Iter(), f)) }

// Fold is a terminal operation that reduces s to a single value by repeatedly
// applying f(acc, key, value), starting from init. Returns init unchanged when
// the sequence is empty. Requires Go 1.27 method-level generics because the
// accumulator type is independent of K and V.
//
//	Of2(xiter.Enumerate(xiter.Range1(5))).Fold(0, func(acc, i, n int) int {
//	    return acc + n
//	})  // returns 10
func (s Seq2[K, V]) Fold[A any](init A, f func(A, K, V) A) A {
	return xiter.Fold2(s.Iter(), init, f)
}

// TryFold is like Fold but stops when f returns an error. Returns the
// accumulator built so far together with the first error; on success returns
// (final accumulator, nil). Requires Go 1.27 method-level generics.
func (s Seq2[K, V]) TryFold[A any](init A, f func(A, K, V) (A, error)) (A, error) {
	return xiter.TryFold2(s.Iter(), init, f)
}
