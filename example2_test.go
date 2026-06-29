package xiter_test

import (
	"fmt"
	"slices"

	"github.com/go-board/xiter"
)

// ============================================================================
// Source
// ============================================================================

func ExampleFromFunc2() {
	i := 0
	seq := xiter.FromFunc2(func() (string, int, bool) {
		i++
		if i > 2 {
			return "", 0, false
		}
		return fmt.Sprintf("k%d", i), i * 10, true
	})
	for k, v := range seq {
		fmt.Printf("%s=%d\n", k, v)
	}
	// Output:
	// k1=10
	// k2=20
}

func ExampleIterate2() {
	// Index/value pairs where value doubles each step, stopping at index 3.
	for k, v := range xiter.Iterate2(0, 1, func(k, v int) (int, int, bool) {
		if k >= 3 {
			return 0, 0, false
		}
		return k + 1, v * 2, true
	}) {
		fmt.Printf("%d=%d\n", k, v)
	}
	// Output:
	// 0=1
	// 1=2
	// 2=4
	// 3=8
}

func ExampleOnce2() {
	for k, v := range xiter.Once2("name", "Alice") {
		fmt.Printf("%s=%s\n", k, v)
	}
	// Output:
	// name=Alice
}

func ExampleEmpty2() {
	n := 0
	for range xiter.Empty2[string, int]() {
		n++
	}
	fmt.Println("count:", n)
	// Output:
	// count: 0
}

func ExampleRepeat2() {
	for k, v := range xiter.Take2(xiter.Repeat2("k", 1), 3) {
		fmt.Printf("%s=%d\n", k, v)
	}
	// Output:
	// k=1
	// k=1
	// k=1
}

// ============================================================================
// Transform
// ============================================================================

func ExampleMap2() {
	seq := xiter.Map2(slices.All([]int{1, 2, 3}), func(k, v int) (string, int) {
		return fmt.Sprintf("k%d", k), v * v
	})
	for k, v := range seq {
		fmt.Printf("%s=%d\n", k, v)
	}
	// Output:
	// k0=1
	// k1=4
	// k2=9
}

func ExampleMapWhile2() {
	seq := xiter.MapWhile2(slices.All([]int{1, 2, 3, 4}), func(k, v int) (int, int, bool) {
		if v < 3 {
			return k, v * 10, true
		}
		return 0, 0, false
	})
	for k, v := range seq {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:10
	// 1:20
}

func ExampleInspect2() {
	seq := xiter.Inspect2(slices.All([]int{1, 2}), func(k, v int) {
		fmt.Printf("inspect: %d=%d\n", k, v)
	})
	for k, v := range seq {
		fmt.Printf("yield: %d=%d\n", k, v)
	}
	// Output:
	// inspect: 0=1
	// yield: 0=1
	// inspect: 1=2
	// yield: 1=2
}

func ExampleJoin() {
	seq := xiter.Join(slices.All([]string{"a", "b"}), func(k int, v string) string {
		return fmt.Sprintf("%d:%s", k, v)
	})
	for v := range seq {
		fmt.Println(v)
	}
	// Output:
	// 0:a
	// 1:b
}

func ExampleKeys() {
	for k := range xiter.Keys(slices.All([]string{"a", "b"})) {
		fmt.Println(k)
	}
	// Output:
	// 0
	// 1
}

func ExampleValues() {
	for v := range xiter.Values(slices.All([]string{"a", "b"})) {
		fmt.Println(v)
	}
	// Output:
	// a
	// b
}

func ExampleSwap() {
	for v, k := range xiter.Swap(slices.All([]string{"a", "b"})) {
		fmt.Printf("%s:%d\n", v, k)
	}
	// Output:
	// a:0
	// b:1
}

// ============================================================================
// Filter / Slice
// ============================================================================

func ExampleFilter2() {
	seq := xiter.Filter2(slices.All([]int{1, 2, 3, 4}), func(k, v int) bool {
		return v%2 == 0
	})
	for k, v := range seq {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 1:2
	// 3:4
}

func ExampleFilterMap2() {
	seq := xiter.FilterMap2(slices.All([]int{1, 2, 3, 4}), func(k, v int) (string, int, bool) {
		if v%2 == 0 {
			return fmt.Sprintf("k%d", k), v * 10, true
		}
		return "", 0, false
	})
	for k, v := range seq {
		fmt.Printf("%s=%d\n", k, v)
	}
	// Output:
	// k1=20
	// k3=40
}

func ExampleTake2() {
	for k, v := range xiter.Take2(slices.All([]int{1, 2, 3, 4}), 2) {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:1
	// 1:2
}

func ExampleTakeWhile2() {
	seq := xiter.TakeWhile2(slices.All([]int{1, 2, 3, 4}), func(k, v int) bool {
		return v < 3
	})
	for k, v := range seq {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:1
	// 1:2
}

func ExampleSkip2() {
	for k, v := range xiter.Skip2(slices.All([]int{1, 2, 3, 4}), 2) {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 2:3
	// 3:4
}

func ExampleSkipWhile2() {
	seq := xiter.SkipWhile2(slices.All([]int{1, 2, 3, 4}), func(k, v int) bool {
		return v < 3
	})
	for k, v := range seq {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 2:3
	// 3:4
}

func ExampleStepBy2() {
	seq := xiter.StepBy2(slices.All([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}), 3)
	for k, v := range seq {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:0
	// 3:3
	// 6:6
	// 9:9
}

func ExampleChain2() {
	seq1 := slices.All([]int{1, 2})
	seq2 := slices.All([]int{3, 4})
	for k, v := range xiter.Chain2(seq1, seq2) {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:1
	// 1:2
	// 0:3
	// 1:4
}

// ============================================================================
// Terminal
// ============================================================================

func ExampleForEach2() {
	xiter.ForEach2(slices.All([]int{1, 2}), func(k, v int) {
		fmt.Printf("%d:%d\n", k, v)
	})
	// Output:
	// 0:1
	// 1:2
}

func ExampleTryForEach2() {
	err := xiter.TryForEach2(slices.All([]int{1, 2, 3}), func(k, v int) error {
		if v == 2 {
			return fmt.Errorf("stop at %d", v)
		}
		fmt.Printf("%d:%d\n", k, v)
		return nil
	})
	fmt.Println("err:", err)
	// Output:
	// 0:1
	// err: stop at 2
}

func ExampleFold2() {
	sum := xiter.Fold2(slices.All([]int{1, 2, 3, 4}), 0, func(acc, k, v int) int {
		return acc + v
	})
	fmt.Println(sum)
	// Output:
	// 10
}

func ExampleReduce2() {
	sumK, sumV, ok := xiter.Reduce2(slices.All([]int{1, 2, 3}), func(ak, av, k, v int) (int, int) {
		return ak + k, av + v
	})
	fmt.Printf("%d,%d,%t\n", sumK, sumV, ok)
	// Output:
	// 3,6,true
}

func ExampleTryReduce2() {
	sumK, sumV, ok, err := xiter.TryReduce2(slices.All([]int{1, 2, 3}), func(ak, av, k, v int) (int, int, error) {
		return ak + k, av + v, nil
	})
	fmt.Printf("%d,%d,%t,%v\n", sumK, sumV, ok, err)
	// Output:
	// 3,6,true,<nil>
}

func ExampleTryFold2() {
	sum, err := xiter.TryFold2(slices.All([]int{1, 2, 3, 4}), 0, func(acc, k, v int) (int, error) {
		return acc + v, nil
	})
	fmt.Printf("%d,%v\n", sum, err)
	// Output:
	// 10,<nil>
}

func ExampleSize2() {
	fmt.Println(xiter.Size2(slices.All([]int{1, 2, 3, 4})))
	// Output:
	// 4
}

func ExampleSizeFunc2() {
	n := xiter.SizeFunc2(slices.All([]int{1, 2, 3, 4}), func(k, v int) bool { return v%2 == 0 })
	fmt.Println(n)
	// Output:
	// 2
}

func ExampleSizeValue2() {
	seq := func(yield func(int, string) bool) {
		pairs := []struct {
			k int
			v string
		}{{1, "a"}, {1, "a"}, {1, "b"}, {2, "a"}}
		for _, p := range pairs {
			if !yield(p.k, p.v) {
				return
			}
		}
	}
	fmt.Println(xiter.SizeValue2(seq, 1, "a"))
	// Output:
	// 2
}

// ============================================================================
// Compare / Search
// ============================================================================

func ExampleContains2() {
	fmt.Println(xiter.Contains2(slices.All([]int{1, 2, 3}), 1, 2))
	fmt.Println(xiter.Contains2(slices.All([]int{1, 2, 3}), 0, 5))
	// Output:
	// true
	// false
}

func ExampleContainsFunc2() {
	fmt.Println(xiter.ContainsFunc2(slices.All([]int{1, 2, 3}), func(k, v int) bool { return v > 2 }))
	fmt.Println(xiter.ContainsFunc2(slices.All([]int{1, 2, 3}), func(k, v int) bool { return v > 10 }))
	// Output:
	// true
	// false
}

func ExampleAny2() {
	fmt.Println(xiter.Any2(slices.All([]int{1, 2, 3}), func(k, v int) bool { return v > 2 }))
	fmt.Println(xiter.Any2(slices.All([]int{1, 2, 3}), func(k, v int) bool { return v > 10 }))
	// Output:
	// true
	// false
}

func ExampleAll2() {
	fmt.Println(xiter.All2(slices.All([]int{1, 2, 3}), func(k, v int) bool { return v > 0 }))
	fmt.Println(xiter.All2(slices.All([]int{1, 2, 3}), func(k, v int) bool { return v < 3 }))
	// Output:
	// true
	// false
}

func ExampleFirst2() {
	k, v, ok := xiter.First2(slices.All([]int{1, 2, 3}))
	fmt.Printf("%d,%d,%t\n", k, v, ok)
	// Output:
	// 0,1,true
}

func ExampleFirstFunc2() {
	k, v, ok := xiter.FirstFunc2(slices.All([]int{1, 2, 3}), func(k, v int) bool { return v > 1 })
	fmt.Printf("%d,%d,%t\n", k, v, ok)
	// Output:
	// 1,2,true
}

func ExampleLast2() {
	k, v, ok := xiter.Last2(slices.All([]int{1, 2, 3}))
	fmt.Printf("%d,%d,%t\n", k, v, ok)
	// Output:
	// 2,3,true
}

func ExampleLastFunc2() {
	k, v, ok := xiter.LastFunc2(slices.All([]int{1, 2, 3, 4}), func(k, v int) bool { return v%2 == 0 })
	fmt.Printf("%d,%d,%t\n", k, v, ok)
	// Output:
	// 3,4,true
}

func ExamplePosition2() {
	i, ok := xiter.Position2(slices.All([]int{1, 2, 3}), func(k, v int) bool { return v == 2 })
	fmt.Printf("%d,%t\n", i, ok)
	// Output:
	// 1,true
}

func ExampleNth2() {
	k, v, ok := xiter.Nth2(slices.All([]int{10, 20, 30, 40}), 2)
	fmt.Printf("%d:%d,%t\n", k, v, ok)
	// Output:
	// 2:30,true
}

func ExampleFindMap2() {
	k, v, ok := xiter.FindMap2(slices.All([]int{1, 2, 3, 4}), func(_, v int) (string, int, bool) {
		if v%2 == 0 {
			return "even", v * 10, true
		}
		return "", 0, false
	})
	fmt.Printf("%s:%d,%t\n", k, v, ok)
	// Output:
	// even:20,true
}

func ExampleCompare2() {
	fmt.Println(xiter.Compare2(slices.All([]int{1, 2}), slices.All([]int{1, 2})))
	fmt.Println(xiter.Compare2(slices.All([]int{1, 2}), slices.All([]int{1, 3})))
	fmt.Println(xiter.Compare2(slices.All([]int{1, 2}), slices.All([]int{1, 2, 3})))
	// Output:
	// 0
	// -1
	// -1
}

func ExampleCompareFunc2() {
	cmp := func(k1, v1, k2, v2 int) int { return v1 - v2 }
	fmt.Println(xiter.CompareFunc2(slices.All([]int{1, 2}), slices.All([]int{1, 2}), cmp))
	fmt.Println(xiter.CompareFunc2(slices.All([]int{1, 2}), slices.All([]int{1, 3}), cmp))
	// Output:
	// 0
	// -1
}

func ExampleEqual2() {
	fmt.Println(xiter.Equal2(slices.All([]int{1, 2}), slices.All([]int{1, 2})))
	fmt.Println(xiter.Equal2(slices.All([]int{1, 2}), slices.All([]int{1, 3})))
	// Output:
	// true
	// false
}

func ExampleEqualFunc2() {
	eq := func(k1, v1, k2, v2 int) bool { return v1 == v2 }
	fmt.Println(xiter.EqualFunc2(slices.All([]int{1, 2}), slices.All([]int{1, 2}), eq))
	fmt.Println(xiter.EqualFunc2(slices.All([]int{1, 2}), slices.All([]int{1, 3}), eq))
	// Output:
	// true
	// false
}
