package collector_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/go-board/xiter"
	"github.com/go-board/xiter/collector"
)

// seqOf returns an iter.Seq[E] yielding the given elements, honoring the
// stop signal from yield.
func seqOf[E any](es ...E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, e := range es {
			if !yield(e) {
				return
			}
		}
	}
}

// pairsOf returns an iter.Seq2[K, V] pairing keys and values element-wise.
func pairsOf[K, V any](keys []K, values []V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := range keys {
			if !yield(keys[i], values[i]) {
				return
			}
		}
	}
}

// printSorted prints a map[K]V in ascending key order so the example output is
// deterministic.
func printSorted[K cmpOrdered, V any](m map[K]V) {
	for _, k := range slices.Sorted(maps.Keys(m)) {
		fmt.Printf("%v:%v ", k, m[k])
	}
	fmt.Println()
}

type cmpOrdered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

func ExampleCollect() {
	s := xiter.Range1(5) // yields 0, 1, 2, 3, 4
	got := collector.Collect(s, collector.ToSlice[int]())
	fmt.Println(got)
	// Output: [0 1 2 3 4]
}

func ExampleCollect2() {
	pairs := pairsOf([]int{1, 2}, []string{"a", "b"})
	got := collector.Collect2(pairs, collector.ToMap2[int, string]())
	printSorted(got)
	// Output:
	// 1:a 2:b
}

func ExampleToSlice() {
	got := collector.ToSlice[int]()(seqOf(1, 2, 3))
	fmt.Println(got)
	// Output: [1 2 3]
}

func ExampleToSet() {
	got := collector.ToSet[int]()(seqOf(1, 2, 2, 3, 1))
	for _, k := range slices.Sorted(maps.Keys(got)) {
		fmt.Println(k)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleToMap() {
	got := collector.ToMap(func(e int) (int, int) { return e, e * e })(seqOf(1, 2, 3))
	printSorted(got)
	// Output:
	// 1:1 2:4 3:9
}

func ExampleToMapMerge() {
	// Keys are e%2; values are summed on collision.
	got := collector.ToMapMerge(
		func(e int) (int, int) { return e % 2, e },
		func(a, b int) int { return a + b },
	)(seqOf(1, 3, 2, 4))
	printSorted(got)
	// Output:
	// 0:6 1:4
}

func ExampleJoining() {
	got := collector.Joining(", ")(seqOf("a", "b", "c"))
	fmt.Println(got)
	// Output: a, b, c
}

func ExampleGroupingBy() {
	got := collector.GroupingBy(func(e int) int { return e % 2 })(seqOf(1, 2, 3, 4, 5))
	printSorted(got)
	// Output:
	// 0:[2 4] 1:[1 3 5]
}

func ExampleGroupingByDownstream() {
	// Count elements per group by wrapping xiter.Size as a Collector.
	counting := collector.Collector[int, int](func(s iter.Seq[int]) int { return xiter.Size(s) })
	got := collector.GroupingByDownstream(
		func(e int) int { return e % 2 },
		counting,
	)(seqOf(1, 2, 3, 4, 5))
	printSorted(got)
	// Output:
	// 0:2 1:3
}

func ExamplePartitioningBy() {
	got := collector.PartitioningBy(func(e int) bool { return e%2 == 0 })(seqOf(1, 2, 3, 4, 5))
	fmt.Printf("Pass:%v Fail:%v\n", got.Pass, got.Fail)
	// Output: Pass:[2 4] Fail:[1 3 5]
}

func ExampleToMap2() {
	pairs := pairsOf([]int{1, 2}, []string{"a", "b"})
	got := collector.ToMap2[int, string]()(pairs)
	printSorted(got)
	// Output:
	// 1:a 2:b
}

func ExampleToMap2Merge() {
	pairs := pairsOf([]int{1, 1, 2}, []int{10, 20, 30})
	got := collector.ToMap2Merge[int, int](func(a, b int) int { return a + b })(pairs)
	printSorted(got)
	// Output:
	// 1:30 2:30
}

func ExampleToKeys() {
	pairs := pairsOf([]int{1, 2, 3}, []string{"a", "b", "c"})
	got := collector.ToKeys[int, string]()(pairs)
	fmt.Println(got)
	// Output: [1 2 3]
}

func ExampleToValues() {
	pairs := pairsOf([]int{1, 2, 3}, []string{"a", "b", "c"})
	got := collector.ToValues[int, string]()(pairs)
	fmt.Println(got)
	// Output: [a b c]
}
