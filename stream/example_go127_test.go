//go:build go1.27

package stream_test

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/go-board/xiter"
	"github.com/go-board/xiter/stream"
)

// ============================================================================
// Seq[E] methods requiring Go 1.27 method-level generics (seq_go127.go)
// ============================================================================

func ExampleSeq_Map() {
	s := stream.Of(xiter.Range1(3)).Map(func(n int) string {
		return fmt.Sprintf("n=%d", n)
	})
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// n=0
	// n=1
	// n=2
}

func ExampleSeq_MapWhile() {
	s := stream.Of(xiter.Range1(5)).MapWhile(func(n int) (int, bool) {
		if n >= 3 {
			return 0, false
		}
		return n * 10, true
	})
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 10
	// 20
}

func ExampleSeq_FilterMap() {
	s := stream.Of(xiter.Range1(6)).FilterMap(func(n int) (string, bool) {
		if n%2 == 0 {
			return fmt.Sprintf("even%d", n), true
		}
		return "", false
	})
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// even0
	// even2
	// even4
}

func ExampleSeq_Split() {
	s := stream.Of(xiter.Range1(4)).Split(func(n int) (int, string) {
		return n, fmt.Sprintf("v%d", n)
	})
	for k, v := range s.Iter() {
		fmt.Printf("%d:%s\n", k, v)
	}
	// Output:
	// 0:v0
	// 1:v1
	// 2:v2
	// 3:v3
}

func ExampleSeq_Zip() {
	s := stream.Of(xiter.Range1(3)).Zip(stream.Of(xiter.Range2(10, 13)))
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:10
	// 1:11
	// 2:12
}

func ExampleSeq_ZipWith() {
	s := stream.Of(xiter.Range1(3)).ZipWith(stream.Of(xiter.Range2(10, 13)), func(a, b int) int {
		return a + b
	})
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 10
	// 12
	// 14
}

func ExampleSeq_Fold() {
	sum := stream.Of(xiter.Range1(5)).Fold(0, func(acc, n int) int {
		return acc + n
	})
	fmt.Println(sum)
	// Output: 10
}

func ExampleSeq_TryFold() {
	sum, err := stream.Of(xiter.Range1(5)).TryFold(0, func(acc, n int) (int, error) {
		if n == 3 {
			return acc + n, errors.New("stop")
		}
		return acc + n, nil
	})
	fmt.Println(sum, err)
	// Output: 6 stop
}

func ExampleSeq_Scan() {
	s := stream.Of(xiter.Range1(5)).Scan(0, func(acc, n int) (int, bool) {
		return acc + n, true
	})
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 3
	// 6
	// 10
}

func ExampleSeq_Collect() {
	// Use the standard library slices.Collect as a collector.
	got := stream.Of(xiter.Range1(5)).Collect(slices.Collect)
	fmt.Println(got)

	// Or any function matching func(iter.Seq[E]) R.
	joined := stream.Of(xiter.Range1(5)).Collect(func(s iter.Seq[int]) string {
		var b []byte
		for v := range s {
			b = append(b, fmt.Sprintf("%d", v)...)
		}
		return string(b)
	})
	fmt.Println(joined)
	// Output:
	// [0 1 2 3 4]
	// 01234
}

// ============================================================================
// Seq2[K, V] methods requiring Go 1.27 method-level generics (seq2_go127.go)
// ============================================================================

func ExampleSeq2_Map() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).
		Map(func(k, v int) (string, int) {
			return fmt.Sprintf("k%d", k), v * 10
		})
	for k, v := range s.Iter() {
		fmt.Printf("%s:%d\n", k, v)
	}
	// Output:
	// k0:100
	// k1:110
	// k2:120
}

func ExampleSeq2_MapWhile() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(5))).
		MapWhile(func(k, v int) (int, int, bool) {
			if k >= 3 {
				return 0, 0, false
			}
			return k, v * 10, true
		})
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:0
	// 1:10
	// 2:20
}

func ExampleSeq2_FilterMap() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(6))).
		FilterMap(func(k, v int) (string, int, bool) {
			if v%2 == 0 {
				return fmt.Sprintf("k%d", k), v, true
			}
			return "", 0, false
		})
	for k, v := range s.Iter() {
		fmt.Printf("%s:%d\n", k, v)
	}
	// Output:
	// k0:0
	// k2:2
	// k4:4
}

func ExampleSeq2_Join() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).
		Join(func(k, v int) string {
			return fmt.Sprintf("%d:%d", k, v)
		})
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0:10
	// 1:11
	// 2:12
}

func ExampleSeq2_Fold() {
	sum := stream.Of2(xiter.Enumerate(xiter.Range1(5))).
		Fold(0, func(acc, k, v int) int {
			return acc + k + v
		})
	fmt.Println(sum)
	// Output: 20
}

func ExampleSeq2_TryFold() {
	sum, err := stream.Of2(xiter.Enumerate(xiter.Range1(5))).
		TryFold(0, func(acc, k, v int) (int, error) {
			if k == 3 {
				return acc + k + v, errors.New("stop")
			}
			return acc + k + v, nil
		})
	fmt.Println(sum, err)
	// Output: 12 stop
}

func ExampleSeq2_Collect() {
	// Use the standard library maps.Collect as a collector.
	m := stream.Of2(xiter.Enumerate(xiter.Range1(3))).Collect(maps.Collect)
	for _, k := range slices.Sorted(maps.Keys(m)) {
		fmt.Printf("%d:%d\n", k, m[k])
	}
	// Output:
	// 0:0
	// 1:1
	// 2:2
}
