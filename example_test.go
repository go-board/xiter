package xiter_test

import (
	"fmt"
	"iter"

	"github.com/go-board/xiter"
)

// ============================================================================
// Source
// ============================================================================

func ExampleRange1() {
	for v := range xiter.Range1(5) {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleRange2() {
	for v := range xiter.Range2(2, 5) {
		fmt.Println(v)
	}
	// Output:
	// 2
	// 3
	// 4
}

func ExampleRange3() {
	for v := range xiter.Range3(1, 10, 2) {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 3
	// 5
	// 7
	// 9
}

func ExampleFromFunc() {
	i := 0
	seq := xiter.FromFunc(func() (int, bool) {
		i++
		if i > 3 {
			return 0, false
		}
		return i, true
	})
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleIterate() {
	// Powers of two, stopping once we reach 16.
	for v := range xiter.Iterate(1, func(x int) (int, bool) {
		if x >= 16 {
			return 0, false
		}
		return x * 2, true
	}) {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 4
	// 8
	// 16
}

func ExampleOnce() {
	for v := range xiter.Once(42) {
		fmt.Println(v)
	}
	// Output:
	// 42
}

func ExampleEmpty() {
	n := 0
	for range xiter.Empty[int]() {
		n++
	}
	fmt.Println("count:", n)
	// Output:
	// count: 0
}

func ExampleRepeat() {
	for v := range xiter.Take(xiter.Repeat(7), 3) {
		fmt.Println(v)
	}
	// Output:
	// 7
	// 7
	// 7
}

// ============================================================================
// Transform
// ============================================================================

func ExampleMap() {
	seq := xiter.Map(xiter.Range1(3), func(x int) int { return x * 2 })
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 2
	// 4
}

func ExampleMapWhile() {
	seq := xiter.MapWhile(xiter.Range1(10), func(x int) (int, bool) {
		if x < 3 {
			return x * 10, true
		}
		return 0, false
	})
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 10
	// 20
}

func ExampleFlatMap() {
	seq := xiter.FlatMap(xiter.Range1(3), func(x int) iter.Seq[int] {
		return xiter.Range1(x + 1)
	})
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 0
	// 1
	// 0
	// 1
	// 2
}

func ExampleFlatten() {
	seqOfSeqs := func(yield func(iter.Seq[int]) bool) {
		if !yield(xiter.Range1(2)) {
			return
		}
		if !yield(xiter.Range2(10, 12)) {
			return
		}
	}
	for v := range xiter.Flatten(seqOfSeqs) {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 10
	// 11
}

func ExampleInspect() {
	seq := xiter.Inspect(xiter.Range1(3), func(v int) {
		fmt.Printf("inspect: %d\n", v)
	})
	for v := range seq {
		fmt.Printf("yield: %d\n", v)
	}
	// Output:
	// inspect: 0
	// yield: 0
	// inspect: 1
	// yield: 1
	// inspect: 2
	// yield: 2
}

func ExampleEnumerate() {
	src := func(yield func(string) bool) {
		for _, s := range []string{"a", "b"} {
			if !yield(s) {
				return
			}
		}
	}
	for i, v := range xiter.Enumerate(src) {
		fmt.Printf("%d:%s\n", i, v)
	}
	// Output:
	// 0:a
	// 1:b
}

func ExampleSplit() {
	src := func(yield func(int) bool) {
		for _, v := range []int{1, 2} {
			if !yield(v) {
				return
			}
		}
	}
	seq := xiter.Split(src, func(e int) (int, int) { return e, e * 10 })
	for k, v := range seq {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 1:10
	// 2:20
}

func ExampleCast() {
	src := func(yield func(any) bool) {
		for _, v := range []any{1, "x", 3} {
			if !yield(v) {
				return
			}
		}
	}
	for v, ok := range xiter.Cast[int](src) {
		fmt.Printf("%d,%t\n", v, ok)
	}
	// Output:
	// 1,true
	// 0,false
	// 3,true
}

func ExampleScan() {
	seq := xiter.Scan(xiter.Range1(5), 0, func(acc, e int) (int, bool) {
		return acc + e, true
	})
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 3
	// 6
	// 10
}

// ============================================================================
// Filter / Slice
// ============================================================================

func ExampleFilter() {
	seq := xiter.Filter(xiter.Range1(10), func(v int) bool { return v%2 == 0 })
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 2
	// 4
	// 6
	// 8
}

func ExampleFilterMap() {
	seq := xiter.FilterMap(xiter.Range1(5), func(x int) (int, bool) {
		if x%2 == 0 {
			return x * 2, true
		}
		return 0, false
	})
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 4
	// 8
}

func ExampleTake() {
	for v := range xiter.Take(xiter.Range1(10), 3) {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 2
}

func ExampleTakeWhile() {
	seq := xiter.TakeWhile(xiter.Range1(10), func(v int) bool { return v < 3 })
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 2
}

func ExampleSkip() {
	for v := range xiter.Skip(xiter.Range1(5), 2) {
		fmt.Println(v)
	}
	// Output:
	// 2
	// 3
	// 4
}

func ExampleSkipWhile() {
	seq := xiter.SkipWhile(xiter.Range1(5), func(v int) bool { return v < 3 })
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 3
	// 4
}

func ExampleStepBy() {
	for v := range xiter.StepBy(xiter.Range1(10), 3) {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 3
	// 6
	// 9
}

func ExampleChain() {
	for v := range xiter.Chain(xiter.Range1(2), xiter.Range2(10, 12)) {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 10
	// 11
}

func ExampleZip() {
	for a, b := range xiter.Zip(xiter.Range1(5), xiter.Range2(10, 13)) {
		fmt.Printf("%d,%d\n", a, b)
	}
	// Output:
	// 0,10
	// 1,11
	// 2,12
}

func ExampleZipWith() {
	seq := xiter.ZipWith(xiter.Range1(3), xiter.Range2(10, 13), func(a, b int) int {
		return a + b
	})
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 10
	// 12
	// 14
}

// ============================================================================
// Terminal
// ============================================================================

func ExampleForEach() {
	xiter.ForEach(xiter.Range1(3), func(v int) {
		fmt.Println(v)
	})
	// Output:
	// 0
	// 1
	// 2
}

func ExampleTryForEach() {
	err := xiter.TryForEach(xiter.Range1(5), func(v int) error {
		if v == 2 {
			return fmt.Errorf("stop at %d", v)
		}
		fmt.Println(v)
		return nil
	})
	fmt.Println("err:", err)
	// Output:
	// 0
	// 1
	// err: stop at 2
}

func ExampleFold() {
	sum := xiter.Fold(xiter.Range1(5), 0, func(acc, e int) int { return acc + e })
	fmt.Println(sum)
	// Output:
	// 10
}

func ExampleReduce() {
	sum, ok := xiter.Reduce(xiter.Range1(5), func(a, b int) int { return a + b })
	fmt.Printf("%d,%t\n", sum, ok)
	// Output:
	// 10,true
}

func ExampleTryReduce() {
	sum, ok, err := xiter.TryReduce(xiter.Range1(5), func(a, b int) (int, error) {
		return a + b, nil
	})
	fmt.Printf("%d,%t,%v\n", sum, ok, err)
	// Output:
	// 10,true,<nil>
}

func ExampleTryFold() {
	sum, err := xiter.TryFold(xiter.Range1(5), 0, func(acc, e int) (int, error) {
		return acc + e, nil
	})
	fmt.Printf("%d,%v\n", sum, err)
	// Output:
	// 10,<nil>
}

func ExampleSize() {
	fmt.Println(xiter.Size(xiter.Range1(5)))
	// Output:
	// 5
}

func ExampleSizeFunc() {
	n := xiter.SizeFunc(xiter.Range1(10), func(v int) bool { return v%2 == 0 })
	fmt.Println(n)
	// Output:
	// 5
}

func ExampleSizeValue() {
	src := func(yield func(int) bool) {
		for _, v := range []int{1, 2, 2, 3, 2} {
			if !yield(v) {
				return
			}
		}
	}
	fmt.Println(xiter.SizeValue(src, 2))
	// Output:
	// 3
}

// ============================================================================
// Compare / Search
// ============================================================================

func ExampleContains() {
	fmt.Println(xiter.Contains(xiter.Range1(5), 3))
	fmt.Println(xiter.Contains(xiter.Range1(5), 10))
	// Output:
	// true
	// false
}

func ExampleContainsFunc() {
	fmt.Println(xiter.ContainsFunc(xiter.Range1(5), func(v int) bool { return v > 3 }))
	fmt.Println(xiter.ContainsFunc(xiter.Range1(5), func(v int) bool { return v > 10 }))
	// Output:
	// true
	// false
}

func ExampleAny() {
	fmt.Println(xiter.Any(xiter.Range1(5), func(v int) bool { return v > 3 }))
	fmt.Println(xiter.Any(xiter.Range1(5), func(v int) bool { return v > 10 }))
	// Output:
	// true
	// false
}

func ExampleAll() {
	fmt.Println(xiter.All(xiter.Range1(5), func(v int) bool { return v >= 0 }))
	fmt.Println(xiter.All(xiter.Range1(5), func(v int) bool { return v < 3 }))
	// Output:
	// true
	// false
}

func ExampleFirst() {
	v, ok := xiter.First(xiter.Range1(5))
	fmt.Printf("%d,%t\n", v, ok)
	// Output:
	// 0,true
}

func ExampleFirstFunc() {
	v, ok := xiter.FirstFunc(xiter.Range1(5), func(v int) bool { return v > 2 })
	fmt.Printf("%d,%t\n", v, ok)
	// Output:
	// 3,true
}

func ExampleLast() {
	v, ok := xiter.Last(xiter.Range1(5))
	fmt.Printf("%d,%t\n", v, ok)
	// Output:
	// 4,true
}

func ExampleLastFunc() {
	v, ok := xiter.LastFunc(xiter.Range1(10), func(v int) bool { return v%2 == 0 })
	fmt.Printf("%d,%t\n", v, ok)
	// Output:
	// 8,true
}

func ExamplePosition() {
	i, ok := xiter.Position(xiter.Range1(5), func(v int) bool { return v == 3 })
	fmt.Printf("%d,%t\n", i, ok)
	// Output:
	// 3,true
}

func ExampleNth() {
	v, ok := xiter.Nth(xiter.Range1(10), 3)
	fmt.Printf("%d,%t\n", v, ok)
	// Output:
	// 3,true
}

func ExampleFindMap() {
	v, ok := xiter.FindMap(xiter.Range1(5), func(n int) (string, bool) {
		if n%2 == 0 {
			return fmt.Sprintf("even%d", n), true
		}
		return "", false
	})
	fmt.Printf("%s,%t\n", v, ok)
	// Output:
	// even0,true
}

func ExampleCompare() {
	fmt.Println(xiter.Compare(xiter.Range1(3), xiter.Range1(3)))
	fmt.Println(xiter.Compare(xiter.Range1(3), xiter.Range1(5)))
	fmt.Println(xiter.Compare(xiter.Range1(5), xiter.Range1(3)))
	// Output:
	// 0
	// -1
	// 1
}

func ExampleCompareFunc() {
	cmp := func(a, b int) int { return a - b }
	fmt.Println(xiter.CompareFunc(xiter.Range1(3), xiter.Range1(3), cmp))
	fmt.Println(xiter.CompareFunc(xiter.Range1(3), xiter.Range2(2, 5), cmp))
	// Output:
	// 0
	// -2
}

func ExampleEqual() {
	fmt.Println(xiter.Equal(xiter.Range1(3), xiter.Range1(3)))
	fmt.Println(xiter.Equal(xiter.Range1(3), xiter.Range1(5)))
	// Output:
	// true
	// false
}

func ExampleEqualFunc() {
	eq := func(a, b int) bool { return a == b }
	fmt.Println(xiter.EqualFunc(xiter.Range1(3), xiter.Range1(3), eq))
	fmt.Println(xiter.EqualFunc(xiter.Range1(3), xiter.Range1(5), eq))
	// Output:
	// true
	// false
}

func ExampleMax() {
	v, ok := xiter.Max(xiter.Range1(5))
	fmt.Printf("%d,%t\n", v, ok)
	// Output:
	// 4,true
}

func ExampleMaxFunc() {
	cmp := func(a, b int) int { return a - b }
	v, ok := xiter.MaxFunc(xiter.Range1(5), cmp)
	fmt.Printf("%d,%t\n", v, ok)
	// Output:
	// 4,true
}

func ExampleMin() {
	v, ok := xiter.Min(xiter.Range1(5))
	fmt.Printf("%d,%t\n", v, ok)
	// Output:
	// 0,true
}

func ExampleMinFunc() {
	cmp := func(a, b int) int { return a - b }
	v, ok := xiter.MinFunc(xiter.Range1(5), cmp)
	fmt.Printf("%d,%t\n", v, ok)
	// Output:
	// 0,true
}

func ExampleMinMax() {
	min, max, ok := xiter.MinMax(xiter.Range1(5))
	fmt.Printf("%d,%d,%t\n", min, max, ok)
	// Output:
	// 0,4,true
}

func ExampleMinMaxFunc() {
	cmp := func(a, b int) int { return a - b }
	min, max, ok := xiter.MinMaxFunc(xiter.Range1(5), cmp)
	fmt.Printf("%d,%d,%t\n", min, max, ok)
	// Output:
	// 0,4,true
}

func ExampleIsSorted() {
	fmt.Println(xiter.IsSorted(xiter.Range1(5)))
	unsorted := func(yield func(int) bool) {
		for _, v := range []int{3, 1, 2} {
			if !yield(v) {
				return
			}
		}
	}
	fmt.Println(xiter.IsSorted(unsorted))
	// Output:
	// true
	// false
}

func ExampleIsSortedFunc() {
	cmp := func(a, b int) int { return a - b }
	ascSrc := func(yield func(int) bool) {
		for _, v := range []int{1, 2, 3} {
			if !yield(v) {
				return
			}
		}
	}
	descSrc := func(yield func(int) bool) {
		for _, v := range []int{3, 2, 1} {
			if !yield(v) {
				return
			}
		}
	}
	fmt.Println(xiter.IsSortedFunc(ascSrc, cmp))
	fmt.Println(xiter.IsSortedFunc(descSrc, cmp))
	// Output:
	// true
	// true
}
