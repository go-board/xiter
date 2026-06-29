package stream_test

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/go-board/xiter"
	"github.com/go-board/xiter/stream"
)

// ============================================================================
// Source constructors (seq.go)
// ============================================================================

func ExampleOf() {
	s := stream.Of(xiter.Range1(3))
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 2
}

func ExampleFromFunc() {
	i := 0
	s := stream.FromFunc(func() (int, bool) {
		if i >= 3 {
			return 0, false
		}
		i++
		return i, true
	})
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleIterate() {
	// Powers of two, stopping once we reach 16.
	s := stream.Iterate(1, func(x int) (int, bool) {
		if x >= 16 {
			return 0, false
		}
		return x * 2, true
	})
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 4
	// 8
	// 16
}

// ============================================================================
// Seq[E] methods (seq.go)
// ============================================================================

func ExampleSeq_Iter() {
	s := stream.Of(xiter.Range1(3)).Take(2)
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
}

func ExampleSeq_Filter() {
	s := stream.Of(xiter.Range1(6)).Filter(func(n int) bool { return n%2 == 0 })
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 2
	// 4
}

func ExampleSeq_Inspect() {
	s := stream.Of(xiter.Range1(3)).
		Inspect(func(n int) { fmt.Println("inspect", n) }).
		Take(2)
	for v := range s.Iter() {
		fmt.Println("emit", v)
	}
	// Output:
	// inspect 0
	// emit 0
	// inspect 1
	// emit 1
}

func ExampleSeq_Take() {
	s := stream.Of(xiter.Range1(5)).Take(2)
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
}

func ExampleSeq_Skip() {
	s := stream.Of(xiter.Range1(5)).Skip(3)
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 3
	// 4
}

func ExampleSeq_TakeWhile() {
	s := stream.Of(xiter.Range1(5)).TakeWhile(func(n int) bool { return n < 3 })
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 2
}

func ExampleSeq_SkipWhile() {
	s := stream.Of(xiter.Range1(5)).SkipWhile(func(n int) bool { return n < 3 })
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 3
	// 4
}

func ExampleSeq_StepBy() {
	s := stream.Of(xiter.Range1(10)).StepBy(3)
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 3
	// 6
	// 9
}

func ExampleSeq_Chain() {
	s := stream.Of(xiter.Range2(0, 2)).Chain(stream.Of(xiter.Range2(10, 12)))
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 10
	// 11
}

func ExampleSeq_Enumerate() {
	s := stream.Of(xiter.Range2(10, 13)).Enumerate()
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:10
	// 1:11
	// 2:12
}

func ExampleSeq_ForEach() {
	stream.Of(xiter.Range1(3)).ForEach(func(n int) {
		fmt.Println(n)
	})
	// Output:
	// 0
	// 1
	// 2
}

func ExampleSeq_TryForEach() {
	err := stream.Of(xiter.Range1(5)).TryForEach(func(n int) error {
		if n == 2 {
			return errors.New("stop")
		}
		fmt.Println(n)
		return nil
	})
	fmt.Println("err:", err)
	// Output:
	// 0
	// 1
	// err: stop
}

func ExampleSeq_Reduce() {
	v, ok := stream.Of(xiter.Range1(5)).Reduce(func(a, b int) int { return a + b })
	fmt.Println(v, ok)
	// Output: 10 true
}

func ExampleSeq_TryReduce() {
	v, ok, err := stream.Of(xiter.Range1(5)).TryReduce(func(acc, n int) (int, error) {
		if n == 3 {
			return acc + n, errors.New("stop")
		}
		return acc + n, nil
	})
	fmt.Println(v, ok, err)
	// Output: 6 true stop
}

func ExampleSeq_Size() {
	fmt.Println(stream.Of(xiter.Range1(5)).Size())
	// Output: 5
}

func ExampleSeq_SizeFunc() {
	n := stream.Of(xiter.Range1(10)).SizeFunc(func(n int) bool { return n%2 == 0 })
	fmt.Println(n)
	// Output: 5
}

func ExampleSeq_Any() {
	fmt.Println(stream.Of(xiter.Range1(5)).Any(func(n int) bool { return n == 3 }))
	// Output: true
}

func ExampleSeq_All() {
	fmt.Println(stream.Of(xiter.Range1(5)).All(func(n int) bool { return n < 10 }))
	// Output: true
}

func ExampleSeq_First() {
	v, ok := stream.Of(xiter.Range2(10, 13)).First()
	fmt.Println(v, ok)
	// Output: 10 true
}

func ExampleSeq_Last() {
	v, ok := stream.Of(xiter.Range2(10, 13)).Last()
	fmt.Println(v, ok)
	// Output: 12 true
}

func ExampleSeq_FirstFunc() {
	v, ok := stream.Of(xiter.Range1(10)).FirstFunc(func(n int) bool { return n > 5 })
	fmt.Println(v, ok)
	// Output: 6 true
}

func ExampleSeq_LastFunc() {
	v, ok := stream.Of(xiter.Range1(10)).LastFunc(func(n int) bool { return n < 5 })
	fmt.Println(v, ok)
	// Output: 4 true
}

func ExampleSeq_Position() {
	i, ok := stream.Of(xiter.Range1(10)).Position(func(n int) bool { return n == 3 })
	fmt.Println(i, ok)
	// Output: 3 true
}

func ExampleSeq_Nth() {
	v, ok := stream.Of(xiter.Range1(10)).Nth(3)
	fmt.Println(v, ok)
	// Output: 3 true
}

func ExampleSeq_IsSortedFunc() {
	fmt.Println(stream.Of(xiter.Range1(5)).IsSortedFunc(cmp.Compare))
	// Output: true
}

func ExampleSeq_CompareFunc() {
	a := stream.Of(xiter.Range1(3))
	b := stream.Of(xiter.Range1(5))
	fmt.Println(a.CompareFunc(b, cmp.Compare))
	// Output: -1
}

func ExampleSeq_EqualFunc() {
	a := stream.Of(xiter.Range1(3))
	b := stream.Of(xiter.Range1(3))
	fmt.Println(a.EqualFunc(b, func(x, y int) bool { return x == y }))
	// Output: true
}

func ExampleSeq_MaxFunc() {
	v, ok := stream.Of(xiter.Range2(1, 6)).MaxFunc(cmp.Compare)
	fmt.Println(v, ok)
	// Output: 5 true
}

func ExampleSeq_MinFunc() {
	v, ok := stream.Of(xiter.Range2(1, 6)).MinFunc(cmp.Compare)
	fmt.Println(v, ok)
	// Output: 1 true
}

func ExampleSeq_MinMaxFunc() {
	min, max, ok := stream.Of(xiter.Range2(1, 6)).MinMaxFunc(cmp.Compare)
	fmt.Println(min, max, ok)
	// Output: 1 5 true
}

func ExampleSeq_ContainsFunc() {
	fmt.Println(stream.Of(xiter.Range1(5)).ContainsFunc(func(n int) bool { return n == 3 }))
	// Output: true
}

// ============================================================================
// Source constructors (seq2.go)
// ============================================================================

func ExampleOf2() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13)))
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:10
	// 1:11
	// 2:12
}

func ExampleFromFunc2() {
	i := 0
	s := stream.FromFunc2(func() (int, string, bool) {
		if i >= 3 {
			return 0, "", false
		}
		i++
		return i, fmt.Sprintf("v%d", i), true
	})
	for k, v := range s.Iter() {
		fmt.Printf("%d:%s\n", k, v)
	}
	// Output:
	// 1:v1
	// 2:v2
	// 3:v3
}

func ExampleIterate2() {
	// Index/value pairs where value doubles each step, stopping at index 3.
	s := stream.Iterate2(0, 1, func(k, v int) (int, int, bool) {
		if k >= 3 {
			return 0, 0, false
		}
		return k + 1, v * 2, true
	})
	for k, v := range s.Iter() {
		fmt.Printf("%d=%d\n", k, v)
	}
	// Output:
	// 0=1
	// 1=2
	// 2=4
	// 3=8
}

// ============================================================================
// Seq2[K, V] methods (seq2.go)
// ============================================================================

func ExampleSeq2_Iter() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).Take(2)
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:10
	// 1:11
}

func ExampleSeq2_Filter() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(6))).
		Filter(func(k, v int) bool { return v%2 == 0 })
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:0
	// 2:2
	// 4:4
}

func ExampleSeq2_Keys() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).Keys()
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 0
	// 1
	// 2
}

func ExampleSeq2_Values() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).Values()
	for v := range s.Iter() {
		fmt.Println(v)
	}
	// Output:
	// 10
	// 11
	// 12
}

func ExampleSeq2_Swap() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).Swap()
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 10:0
	// 11:1
	// 12:2
}

func ExampleSeq2_Inspect() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).
		Inspect(func(k, v int) { fmt.Println("inspect", k, v) }).
		Take(2)
	for k, v := range s.Iter() {
		fmt.Println("emit", k, v)
	}
	// Output:
	// inspect 0 10
	// emit 0 10
	// inspect 1 11
	// emit 1 11
}

func ExampleSeq2_Take() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).Take(2)
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:10
	// 1:11
}

func ExampleSeq2_Skip() {
	s := stream.Of2(xiter.Enumerate(xiter.Range2(10, 15))).Skip(3)
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 3:13
	// 4:14
}

func ExampleSeq2_TakeWhile() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(5))).
		TakeWhile(func(k, v int) bool { return k < 3 })
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:0
	// 1:1
	// 2:2
}

func ExampleSeq2_SkipWhile() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(5))).
		SkipWhile(func(k, v int) bool { return k < 3 })
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 3:3
	// 4:4
}

func ExampleSeq2_StepBy() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(10))).StepBy(3)
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:0
	// 3:3
	// 6:6
	// 9:9
}

func ExampleSeq2_Chain() {
	a := stream.Of2(xiter.Enumerate(xiter.Range2(0, 2)))
	b := stream.Of2(xiter.Enumerate(xiter.Range2(10, 12)))
	s := a.Chain(b)
	for k, v := range s.Iter() {
		fmt.Printf("%d:%d\n", k, v)
	}
	// Output:
	// 0:0
	// 1:1
	// 0:10
	// 1:11
}

func ExampleSeq2_ForEach() {
	stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).ForEach(func(k, v int) {
		fmt.Printf("%d:%d\n", k, v)
	})
	// Output:
	// 0:10
	// 1:11
	// 2:12
}

func ExampleSeq2_TryForEach() {
	err := stream.Of2(xiter.Enumerate(xiter.Range1(5))).TryForEach(func(k, v int) error {
		if k == 2 {
			return errors.New("stop")
		}
		fmt.Printf("%d:%d\n", k, v)
		return nil
	})
	fmt.Println("err:", err)
	// Output:
	// 0:0
	// 1:1
	// err: stop
}

func ExampleSeq2_Reduce() {
	k, v, ok := stream.Of2(xiter.Enumerate(xiter.Range1(5))).Reduce(func(accK, accV, k, v int) (int, int) {
		return accK + k, accV + v
	})
	fmt.Println(k, v, ok)
	// Output: 10 10 true
}

func ExampleSeq2_TryReduce() {
	k, v, ok, err := stream.Of2(xiter.Enumerate(xiter.Range1(5))).TryReduce(func(accK, accV, k, v int) (int, int, error) {
		if k == 3 {
			return accK + k, accV + v, errors.New("stop")
		}
		return accK + k, accV + v, nil
	})
	fmt.Println(k, v, ok, err)
	// Output: 6 6 true stop
}

func ExampleSeq2_Size() {
	fmt.Println(stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).Size())
	// Output: 3
}

func ExampleSeq2_SizeFunc() {
	n := stream.Of2(xiter.Enumerate(xiter.Range1(6))).
		SizeFunc(func(k, v int) bool { return v%2 == 0 })
	fmt.Println(n)
	// Output: 3
}

func ExampleSeq2_Any() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(5)))
	fmt.Println(s.Any(func(k, v int) bool { return v == 3 }))
	// Output: true
}

func ExampleSeq2_All() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(5)))
	fmt.Println(s.All(func(k, v int) bool { return k == v }))
	// Output: true
}

func ExampleSeq2_First() {
	k, v, ok := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).First()
	fmt.Println(k, v, ok)
	// Output: 0 10 true
}

func ExampleSeq2_Last() {
	k, v, ok := stream.Of2(xiter.Enumerate(xiter.Range2(10, 13))).Last()
	fmt.Println(k, v, ok)
	// Output: 2 12 true
}

func ExampleSeq2_FirstFunc() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(10)))
	k, v, ok := s.FirstFunc(func(k, v int) bool { return v > 5 })
	fmt.Println(k, v, ok)
	// Output: 6 6 true
}

func ExampleSeq2_LastFunc() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(10)))
	k, v, ok := s.LastFunc(func(k, v int) bool { return v < 5 })
	fmt.Println(k, v, ok)
	// Output: 4 4 true
}

func ExampleSeq2_Position() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(5)))
	i, ok := s.Position(func(k, v int) bool { return v == 3 })
	fmt.Println(i, ok)
	// Output: 3 true
}

func ExampleSeq2_Nth() {
	k, v, ok := stream.Of2(xiter.Enumerate(xiter.Range1(10))).Nth(3)
	fmt.Printf("%d:%d,%t\n", k, v, ok)
	// Output: 3:3,true
}

func ExampleSeq2_CompareFunc() {
	a := stream.Of2(xiter.Enumerate(xiter.Range1(3)))
	b := stream.Of2(xiter.Enumerate(xiter.Range1(5)))
	fmt.Println(a.CompareFunc(b, func(k1, v1, k2, v2 int) int { return cmp.Compare(v1, v2) }))
	// Output: -1
}

func ExampleSeq2_EqualFunc() {
	a := stream.Of2(xiter.Enumerate(xiter.Range1(3)))
	b := stream.Of2(xiter.Enumerate(xiter.Range1(3)))
	fmt.Println(a.EqualFunc(b, func(k1, v1, k2, v2 int) bool { return k1 == k2 && v1 == v2 }))
	// Output: true
}

func ExampleSeq2_ContainsFunc() {
	s := stream.Of2(xiter.Enumerate(xiter.Range1(5)))
	fmt.Println(s.ContainsFunc(func(k, v int) bool { return v == 3 }))
	// Output: true
}
