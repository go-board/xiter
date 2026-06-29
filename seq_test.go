package xiter

import (
	"errors"
	"iter"
	"reflect"
	"testing"
)

// ============================================================================
// Source
// ============================================================================

func TestRange1(t *testing.T) {
	s := Range1(5)
	result := ToSlice(s)
	expected := []int{0, 1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("Range1(5) returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Range1(5)[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	s = Range1(0)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range1(0) returned %d elements, expected 0", len(result))
	}

	s = Range1(1)
	result = ToSlice(s)
	expected = []int{0}
	if len(result) != len(expected) {
		t.Errorf("Range1(1) returned %d elements, expected %d", len(result), len(expected))
	}
	if result[0] != expected[0] {
		t.Errorf("Range1(1)[0] = %d, expected %d", result[0], expected[0])
	}
}

func TestRange2(t *testing.T) {
	s := Range2(2, 7)
	result := ToSlice(s)
	expected := []int{2, 3, 4, 5, 6}
	if len(result) != len(expected) {
		t.Errorf("Range2(2,7) returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Range2(2,7)[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	s = Range2(5, 5)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range2(5,5) returned %d elements, expected 0", len(result))
	}

	s = Range2(7, 2)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range2(7,2) returned %d elements, expected 0", len(result))
	}
}

func TestRange3(t *testing.T) {
	s := Range3(1, 10, 2)
	result := ToSlice(s)
	expected := []int{1, 3, 5, 7, 9}
	if len(result) != len(expected) {
		t.Errorf("Range3(1,10,2) returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Range3(1,10,2)[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	s = Range3(1, 10, 0)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range3(1,10,0) returned %d elements, expected 0", len(result))
	}

	s = Range3(10, 1, 2)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range3(10,1,2) returned %d elements, expected 0", len(result))
	}
}

func TestFromFunc(t *testing.T) {
	count := 0
	s := FromFunc(func() (int, bool) {
		count++
		if count <= 3 {
			return count, true
		}
		return 0, false
	})
	result := ToSlice(s)
	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("FromFunc returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("FromFunc[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	s = FromFunc(func() (int, bool) {
		return 0, false
	})
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("FromFunc returned %d elements, expected 0", len(result))
	}
}

func TestIterate(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		got := ToSlice(Iterate(1, func(x int) (int, bool) {
			if x >= 16 {
				return 0, false
			}
			return x * 2, true
		}))
		want := []int{1, 2, 4, 8, 16}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i, v := range got {
			if v != want[i] {
				t.Fatalf("got[%d]=%d, want %d", i, v, want[i])
			}
		}
	})
	t.Run("seed only when next immediately stops", func(t *testing.T) {
		got := ToSlice(Iterate(42, func(x int) (int, bool) { return 0, false }))
		if len(got) != 1 || got[0] != 42 {
			t.Fatalf("got %v, want [42]", got)
		}
	})
	t.Run("early stop", func(t *testing.T) {
		// Break after 3 yields; Iterate must honor the stop signal and not
		// continue calling next forever.
		var nextCalls int
		s := Iterate(0, func(x int) (int, bool) {
			nextCalls++
			if x >= 2 {
				return 0, false
			}
			return x + 1, true
		})
		var got []int
		n := 0
		s(func(v int) bool {
			n++
			if n > 3 {
				return false
			}
			got = append(got, v)
			return true
		})
		if len(got) != 3 {
			t.Fatalf("got %d elements, want 3", len(got))
		}
		if nextCalls == 0 {
			t.Fatalf("next was never called")
		}
	})
	t.Run("yield false on seed", func(t *testing.T) {
		// Consumer rejects the very first element (the seed): next must never
		// be called.
		var nextCalls int
		s := Iterate(0, func(x int) (int, bool) {
			nextCalls++
			return x + 1, true
		})
		s(func(v int) bool { return false })
		if nextCalls != 0 {
			t.Fatalf("next was called %d times, want 0", nextCalls)
		}
	})
	t.Run("yield false after first element", func(t *testing.T) {
		// Consumer accepts the seed but rejects the second element: next must
		// not be called again.
		var nextCalls int
		s := Iterate(0, func(x int) (int, bool) {
			nextCalls++
			return x + 1, true
		})
		n := 0
		s(func(v int) bool {
			n++
			return n == 1
		})
		if nextCalls != 1 {
			t.Fatalf("next was called %d times, want 1", nextCalls)
		}
	})
}

func TestOnce(t *testing.T) {
	s := Once(42)
	result := ToSlice(s)
	expected := []int{42}
	if len(result) != 1 {
		t.Errorf("Once returned %d elements, expected 1", len(result))
	}
	if result[0] != expected[0] {
		t.Errorf("Once returned %d, expected %d", result[0], expected[0])
	}
}

func TestEmpty(t *testing.T) {
	s := Empty[int]()
	result := ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Empty returned %d elements, expected 0", len(result))
	}
}

func TestRepeat(t *testing.T) {
	s := Repeat("test")
	taken := Take(s, 3)
	result := ToSlice(taken)
	expected := []string{"test", "test", "test"}
	if len(result) != len(expected) {
		t.Errorf("Repeat returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Repeat[%d] = %s, expected %s", i, v, expected[i])
		}
	}
}

// ============================================================================
// Transform
// ============================================================================

func TestMap(t *testing.T) {
	s := Range1(5)
	mapped := Map(s, func(x int) int { return x * 2 })
	result := ToSlice(mapped)
	expected := []int{0, 2, 4, 6, 8}
	if len(result) != len(expected) {
		t.Errorf("Map returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Map[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	mappedStr := Map(s, func(x int) string { return string(rune('a' + x)) })
	resultStr := ToSlice(mappedStr)
	expectedStr := []string{"a", "b", "c", "d", "e"}
	if len(resultStr) != len(expectedStr) {
		t.Errorf("Map returned %d elements, expected %d", len(resultStr), len(expectedStr))
	}
	for i, v := range resultStr {
		if v != expectedStr[i] {
			t.Errorf("Map[%d] = %s, expected %s", i, v, expectedStr[i])
		}
	}

	emptySeq := Empty[int]()
	mapped = Map(emptySeq, func(x int) int { return x * 2 })
	result = ToSlice(mapped)
	if len(result) != 0 {
		t.Errorf("Map on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestFilterMap(t *testing.T) {
	source := func(yield func(int) bool) {
		for i := 1; i <= 5; i++ {
			if !yield(i) {
				return
			}
		}
	}
	filterFn := func(x int) (int, bool) {
		if x%2 == 0 {
			return x * 2, true
		}
		return 0, false
	}
	result := FilterMap(source, filterFn)
	expected := []int{4, 8}
	var actual []int
	for x := range result {
		actual = append(actual, x)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FilterMap failed. Expected: %v, Got: %v", expected, actual)
	}
}

func TestInspect(t *testing.T) {
	seen := []int{}
	got := ToSlice(Inspect(Range1(5), func(v int) {
		seen = append(seen, v)
	}))
	want := []int{0, 1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if !reflect.DeepEqual(seen, want) {
		t.Fatalf("seen %v, want %v", seen, want)
	}
}

func TestInspectStopsWithConsumer(t *testing.T) {
	seen := []int{}
	got := ToSlice(Take(Inspect(Range1(5), func(v int) {
		seen = append(seen, v)
	}), 2))
	want := []int{0, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if !reflect.DeepEqual(seen, want) {
		t.Fatalf("seen %v, want %v", seen, want)
	}
}

func TestScan(t *testing.T) {
	got := ToSlice(Scan(Range1(5), 0, func(acc, e int) (int, bool) {
		return acc + e, true
	}))
	want := []int{0, 1, 3, 6, 10}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestScanEmpty(t *testing.T) {
	got := ToSlice(Scan(Empty[int](), 0, func(acc, e int) (int, bool) {
		return acc + e, true
	}))
	if len(got) != 0 {
		t.Fatalf("got %v, want empty", got)
	}
}

func TestScanStopEarly(t *testing.T) {
	got := ToSlice(Scan(Range1(10), 0, func(acc, e int) (int, bool) {
		next := acc + e
		if next >= 3 {
			return next, false
		}
		return next, true
	}))
	want := []int{0, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestScanNonZeroInit(t *testing.T) {
	got := ToSlice(Scan(func(yield func(string) bool) {
		for _, s := range []string{"a", "b", "c"} {
			if !yield(s) {
				return
			}
		}
	}, "x", func(acc, e string) (string, bool) {
		return acc + e, true
	}))
	want := []string{"xa", "xab", "xabc"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// ============================================================================
// Filter / Slice
// ============================================================================

func TestTake(t *testing.T) {
	s := Range1(10)
	taken := Take(s, 5)
	result := ToSlice(taken)
	expected := []int{0, 1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("Take returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Take[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	s = Range1(5)
	taken = Take(s, 10)
	result = ToSlice(taken)
	expected = []int{0, 1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("Take returned %d elements, expected %d", len(result), len(expected))
	}

	taken = Take(s, 0)
	result = ToSlice(taken)
	if len(result) != 0 {
		t.Errorf("Take returned %d elements, expected 0", len(result))
	}

	emptySeq := Empty[int]()
	taken = Take(emptySeq, 5)
	result = ToSlice(taken)
	if len(result) != 0 {
		t.Errorf("Take on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestSkip(t *testing.T) {
	s := Range1(10)
	skiped := Skip(s, 5)
	result := ToSlice(skiped)
	expected := []int{5, 6, 7, 8, 9}
	if len(result) != len(expected) {
		t.Errorf("Skip returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Skip[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	s = Range1(5)
	skiped = Skip(s, 10)
	result = ToSlice(skiped)
	if len(result) != 0 {
		t.Errorf("Skip returned %d elements, expected 0", len(result))
	}

	skiped = Skip(s, 0)
	result = ToSlice(skiped)
	expected = []int{0, 1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("Skip returned %d elements, expected %d", len(result), len(expected))
	}

	emptySeq := Empty[int]()
	skiped = Skip(emptySeq, 5)
	result = ToSlice(skiped)
	if len(result) != 0 {
		t.Errorf("Skip on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestFilter(t *testing.T) {
	s := Range1(10)
	filtered := Filter(s, func(x int) bool { return x%2 == 0 })
	result := ToSlice(filtered)
	expected := []int{0, 2, 4, 6, 8}
	if len(result) != len(expected) {
		t.Errorf("Filter returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Filter[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	filtered = Filter(s, func(x int) bool { return true })
	result = ToSlice(filtered)
	expected = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if len(result) != len(expected) {
		t.Errorf("Filter returned %d elements, expected %d", len(result), len(expected))
	}

	filtered = Filter(s, func(x int) bool { return false })
	result = ToSlice(filtered)
	if len(result) != 0 {
		t.Errorf("Filter returned %d elements, expected 0", len(result))
	}

	emptySeq := Empty[int]()
	filtered = Filter(emptySeq, func(x int) bool { return x%2 == 0 })
	result = ToSlice(filtered)
	if len(result) != 0 {
		t.Errorf("Filter on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestTakeWhile(t *testing.T) {
	s := Range1(10)
	takenWhile := TakeWhile(s, func(x int) bool { return x <= 5 })
	result := ToSlice(takenWhile)
	expected := []int{0, 1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("TakeWhile returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("TakeWhile[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	takenWhile = TakeWhile(s, func(x int) bool { return false })
	result = ToSlice(takenWhile)
	if len(result) != 0 {
		t.Errorf("TakeWhile returned %d elements, expected 0", len(result))
	}

	takenWhile = TakeWhile(s, func(x int) bool { return true })
	result = ToSlice(takenWhile)
	expected = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if len(result) != len(expected) {
		t.Errorf("TakeWhile returned %d elements, expected %d", len(result), len(expected))
	}

	emptySeq := Empty[int]()
	takenWhile = TakeWhile(emptySeq, func(x int) bool { return true })
	result = ToSlice(takenWhile)
	if len(result) != 0 {
		t.Errorf("TakeWhile on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestSkipWhile(t *testing.T) {
	s := Range1(10)
	skipedWhile := SkipWhile(s, func(x int) bool { return x <= 5 })
	result := ToSlice(skipedWhile)
	expected := []int{6, 7, 8, 9}
	if len(result) != len(expected) {
		t.Errorf("SkipWhile returned %d elements, expected %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("SkipWhile[%d] = %d, expected %d", i, v, expected[i])
		}
	}

	skipedWhile = SkipWhile(s, func(x int) bool { return false })
	result = ToSlice(skipedWhile)
	expected = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if len(result) != len(expected) {
		t.Errorf("SkipWhile returned %d elements, expected %d", len(result), len(expected))
	}

	skipedWhile = SkipWhile(s, func(x int) bool { return true })
	result = ToSlice(skipedWhile)
	if len(result) != 0 {
		t.Errorf("SkipWhile returned %d elements, expected 0", len(result))
	}

	emptySeq := Empty[int]()
	skipedWhile = SkipWhile(emptySeq, func(x int) bool { return true })
	result = ToSlice(skipedWhile)
	if len(result) != 0 {
		t.Errorf("SkipWhile on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestStepBy(t *testing.T) {
	cases := []struct {
		name string
		src  iter.Seq[int]
		n    int
		want []int
	}{
		{"normal", Range1(10), 3, []int{0, 3, 6, 9}},
		{"step_one", Range1(4), 1, []int{0, 1, 2, 3}},
		{"step_equals_len", Range1(4), 4, []int{0}},
		{"step_greater_than_len", Range1(4), 10, []int{0}},
		{"zero_step", Range1(5), 0, nil},
		{"negative_step", Range1(5), -1, nil},
		{"empty", Empty[int](), 2, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ToSlice(StepBy(c.src, c.n))
			if len(got) != len(c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Fatalf("got %v, want %v", got, c.want)
				}
			}
		})
	}
}

func TestEnumerate(t *testing.T) {
	values := []string{"a", "b", "c", "d", "e"}
	s := func(yield func(string) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
	enumerated := Enumerate(s)
	result := ToMap(enumerated)
	expected := map[int]string{0: "a", 1: "b", 2: "c", 3: "d", 4: "e"}
	if len(result) != len(expected) {
		t.Errorf("Enumerate returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Enumerate[%d] = %s, expected %s", k, result[k], v)
		}
	}

	emptySeq := Empty[string]()
	enumerated = Enumerate(emptySeq)
	result = ToMap(enumerated)
	if len(result) != 0 {
		t.Errorf("Enumerate on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestZip(t *testing.T) {
	got := ToMap(Zip(Range1(5), Range2(10, 13)))
	want := map[int]int{0: 10, 1: 11, 2: 12}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestZipStopsWithConsumer(t *testing.T) {
	seenLeft := []int{}
	seenRight := []int{}
	zipped := Zip(
		Inspect(Range1(5), func(v int) { seenLeft = append(seenLeft, v) }),
		Inspect(Range2(10, 15), func(v int) { seenRight = append(seenRight, v) }),
	)
	got := ToMap(Take2(zipped, 2))
	want := map[int]int{0: 10, 1: 11}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if !reflect.DeepEqual(seenLeft, []int{0, 1}) {
		t.Fatalf("seen left %v, want [0 1]", seenLeft)
	}
	if !reflect.DeepEqual(seenRight, []int{10, 11}) {
		t.Fatalf("seen right %v, want [10 11]", seenRight)
	}
}

func TestZipStopsBeforePullingRightWhenLeftEnds(t *testing.T) {
	seenRight := []int{}
	got := ToMap(Zip(
		Empty[int](),
		Inspect(Range1(3), func(v int) { seenRight = append(seenRight, v) }),
	))
	if len(got) != 0 {
		t.Fatalf("got %v, want empty map", got)
	}
	if len(seenRight) != 0 {
		t.Fatalf("seen right %v, want empty", seenRight)
	}
}

func TestZipWith(t *testing.T) {
	got := ToSlice(ZipWith(Range1(5), Range2(10, 13), func(a, b int) string {
		return string(rune('a'+a)) + string(rune('a'+b-10))
	}))
	want := []string{"aa", "bb", "cc"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// ============================================================================
// Terminal
// ============================================================================

func TestFold(t *testing.T) {
	s := Range1(5)
	sum := Fold(s, 0, func(acc, e int) int { return acc + e })
	expectedSum := 10
	if sum != expectedSum {
		t.Errorf("Fold sum: expected %d, got %d", expectedSum, sum)
	}

	sStr := func(yield func(string) bool) {
		yield("a")
		yield("b")
		yield("c")
	}
	concatenated := Fold(sStr, "", func(acc, e string) string { return acc + e })
	expectedConcatenated := "abc"
	if concatenated != expectedConcatenated {
		t.Errorf("Fold concatenation: expected %s, got %s", expectedConcatenated, concatenated)
	}

	emptySeq := Empty[int]()
	emptyResult := Fold(emptySeq, 42, func(acc, e int) int { return acc + e })
	if emptyResult != 42 {
		t.Errorf("Fold on empty sequence: expected 42, got %d", emptyResult)
	}
}

func TestReduce(t *testing.T) {
	got, ok := Reduce(Range1(5), func(acc, v int) int { return acc + v })
	if !ok {
		t.Fatal("Reduce returned ok=false, want true")
	}
	if got != 10 {
		t.Fatalf("got %d, want 10", got)
	}
}

func TestReduceEmpty(t *testing.T) {
	got, ok := Reduce(Empty[int](), func(acc, v int) int { return acc + v })
	if ok {
		t.Fatal("Reduce returned ok=true, want false")
	}
	if got != 0 {
		t.Fatalf("got %d, want zero", got)
	}
}

func TestTryReduce(t *testing.T) {
	wantErr := errors.New("stop")
	got, ok, err := TryReduce(Range1(5), func(acc, v int) (int, error) {
		acc += v
		if v == 3 {
			return acc, wantErr
		}
		return acc, nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if !ok {
		t.Fatal("TryReduce returned ok=false, want true")
	}
	if got != 6 {
		t.Fatalf("got %d, want 6", got)
	}
}

func TestTryReduceEmpty(t *testing.T) {
	got, ok, err := TryReduce(Empty[int](), func(acc, v int) (int, error) {
		return acc + v, nil
	})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if ok {
		t.Fatal("TryReduce returned ok=true, want false")
	}
	if got != 0 {
		t.Fatalf("got %d, want zero", got)
	}
}

func TestTryForEach(t *testing.T) {
	wantErr := errors.New("stop")
	visited := 0
	err := TryForEach(Range1(5), func(v int) error {
		visited++
		if v == 2 {
			return wantErr
		}
		return nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if visited != 3 {
		t.Fatalf("visited %d elements, want 3", visited)
	}
}

func TestTryFold(t *testing.T) {
	wantErr := errors.New("stop")
	got, err := TryFold(Range1(5), 0, func(acc, v int) (int, error) {
		acc += v
		if v == 3 {
			return acc, wantErr
		}
		return acc, nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if got != 6 {
		t.Fatalf("got accumulator %d, want 6", got)
	}
}

func TestSize(t *testing.T) {
	s := Range1(5)
	size := Size(s)
	expected := 5
	if size != expected {
		t.Errorf("Size: expected %d, got %d", expected, size)
	}

	emptySeq := Empty[int]()
	emptySize := Size(emptySeq)
	if emptySize != 0 {
		t.Errorf("Size on empty sequence: expected 0, got %d", emptySize)
	}

	infinite := Repeat("test")
	limited := Take(infinite, 10)
	limitedSize := Size(limited)
	if limitedSize != 10 {
		t.Errorf("Size on limited infinite sequence: expected 10, got %d", limitedSize)
	}
}

func TestSizeFunc(t *testing.T) {
	s := Range1(10)
	evenCount := SizeFunc(s, func(e int) bool { return e%2 == 0 })
	expected := 5
	if evenCount != expected {
		t.Errorf("SizeFunc (even numbers): expected %d, got %d", expected, evenCount)
	}

	smallerCount := SizeFunc(s, func(e int) bool { return e > 10 })
	if smallerCount != 0 {
		t.Errorf("SizeFunc (numbers > 10): expected 0, got %d", smallerCount)
	}

	emptySeq := Empty[int]()
	emptyCount := SizeFunc(emptySeq, func(e int) bool { return true })
	if emptyCount != 0 {
		t.Errorf("SizeFunc on empty sequence: expected 0, got %d", emptyCount)
	}
}

// ============================================================================
// Compare / Search
// ============================================================================

func TestAny(t *testing.T) {
	numbers := Range1(10)
	result := Any(numbers, func(x int) bool { return x%2 == 0 })
	if !result {
		t.Errorf("Any should return true for even numbers in 0-9")
	}

	numbers = Range1(10)
	result = Any(numbers, func(x int) bool { return x > 10 })
	if result {
		t.Errorf("Any should return false for numbers > 10 in 0-9")
	}

	empty := Empty[int]()
	result = Any(empty, func(x int) bool { return true })
	if result {
		t.Errorf("Any should return false for empty sequence")
	}
}

func TestAll(t *testing.T) {
	numbers := Range2(1, 10)
	result := All(numbers, func(x int) bool { return x > 0 })
	if !result {
		t.Errorf("All should return true for all numbers > 0 in 1-9")
	}

	numbers = Range2(1, 10)
	result = All(numbers, func(x int) bool { return x < 5 })
	if result {
		t.Errorf("All should return false for numbers < 5 in 1-9")
	}

	empty := Empty[int]()
	result = All(empty, func(x int) bool { return true })
	if !result {
		t.Errorf("All should return true for empty sequence")
	}
}

func TestCompare(t *testing.T) {
	s1 := Range1(3)
	s2 := Range1(3)
	result := Compare(s1, s2)
	if result != 0 {
		t.Errorf("Compare (equal sequences): expected 0, got %d", result)
	}

	s3 := Range1(2)
	s4 := Range1(3)
	result = Compare(s3, s4)
	if result >= 0 {
		t.Errorf("Compare (s1 < s2): expected < 0, got %d", result)
	}

	s5 := Range1(3)
	s6 := Range1(2)
	result = Compare(s5, s6)
	if result <= 0 {
		t.Errorf("Compare (s1 > s2): expected > 0, got %d", result)
	}

	s7 := func(yield func(int) bool) {
		yield(1)
		yield(2)
		yield(3)
	}
	s8 := func(yield func(int) bool) {
		yield(1)
		yield(3)
		yield(2)
	}
	result = Compare(s7, s8)
	if result >= 0 {
		t.Errorf("Compare (element values differ): expected < 0, got %d", result)
	}
}

func TestCompareFunc(t *testing.T) {
	s1 := func(yield func(int) bool) {
		yield(1)
		yield(2)
		yield(3)
	}
	s2 := func(yield func(int) bool) {
		yield(1)
		yield(2)
		yield(3)
	}
	result := CompareFunc(s1, s2, func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if result != 0 {
		t.Errorf("CompareFunc (equal sequences): expected 0, got %d", result)
	}

	s3 := func(yield func(int) bool) {
		yield(-1)
		yield(2)
		yield(-3)
	}
	s4 := func(yield func(int) bool) {
		yield(1)
		yield(2)
		yield(3)
	}
	result = CompareFunc(s3, s4, func(a, b int) int {
		absA := a
		if absA < 0 {
			absA = -absA
		}
		absB := b
		if absB < 0 {
			absB = -absB
		}
		if absA < absB {
			return -1
		} else if absA > absB {
			return 1
		}
		return 0
	})
	if result != 0 {
		t.Errorf("CompareFunc (absolute values equal): expected 0, got %d", result)
	}
}

func TestMinMax(t *testing.T) {
	min, max, ok := MinMax(Range1(5))
	if !ok {
		t.Fatal("MinMax returned ok=false, want true")
	}
	if min != 0 || max != 4 {
		t.Fatalf("got (%d, %d), want (0, 4)", min, max)
	}

	min, max, ok = MinMax(Once(7))
	if !ok {
		t.Fatal("MinMax returned ok=false, want true")
	}
	if min != 7 || max != 7 {
		t.Fatalf("got (%d, %d), want (7, 7)", min, max)
	}

	src := func(yield func(int) bool) {
		for _, v := range []int{3, -1, 10, 4, -7, 2} {
			if !yield(v) {
				return
			}
		}
	}
	min, max, ok = MinMax(src)
	if !ok {
		t.Fatal("MinMax returned ok=false, want true")
	}
	if min != -7 || max != 10 {
		t.Fatalf("got (%d, %d), want (-7, 10)", min, max)
	}
}

func TestMinMaxEmpty(t *testing.T) {
	min, max, ok := MinMax(Empty[int]())
	if ok {
		t.Fatal("MinMax returned ok=true, want false")
	}
	if min != 0 || max != 0 {
		t.Fatalf("got (%d, %d), want zero values", min, max)
	}
}

func TestMinMaxFunc(t *testing.T) {
	type person struct {
		name string
		age  int
	}
	src := func(yield func(person) bool) {
		for _, p := range []person{
			{"a", 30},
			{"b", 20},
			{"c", 40},
		} {
			if !yield(p) {
				return
			}
		}
	}
	byAge := func(a, b person) int { return a.age - b.age }
	min, max, ok := MinMaxFunc(src, byAge)
	if !ok {
		t.Fatal("MinMaxFunc returned ok=false, want true")
	}
	if min.age != 20 || max.age != 40 {
		t.Fatalf("got ages (%d, %d), want (20, 40)", min.age, max.age)
	}
}

// ============================================================================
// Supplementary coverage
// ============================================================================

// stopEarly consumes s and stops after the first element, triggering the
// !yield early-return branch in generator-style functions.
func stopEarly[E any](s iter.Seq[E]) { s(func(E) bool { return false }) }
func stopEarly2[K, V any](s iter.Seq2[K, V]) { s(func(K, V) bool { return false }) }

// seqOf builds a Seq from the given elements, honoring yield's stop signal.
func seqOf[E any](vs ...E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, v := range vs {
			if !yield(v) {
				return
			}
		}
	}
}

// seq2Of builds a Seq2 from the given pairs, honoring yield's stop signal.
func seq2Of[K comparable, V any](pairs ...struct {
	K K
	V V
}) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, p := range pairs {
			if !yield(p.K, p.V) {
				return
			}
		}
	}
}

func TestRange3NegativeStep(t *testing.T) {
	got := ToSlice(Range3(10, 1, -2))
	want := []int{10, 8, 6, 4, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got := ToSlice(Range3(1, 10, -2)); len(got) != 0 {
		t.Fatalf("got %v, want empty", got)
	}
}

func TestMapWhile(t *testing.T) {
	got := ToSlice(MapWhile(Range1(10), func(x int) (int, bool) {
		if x < 3 {
			return x * 10, true
		}
		return 0, false
	}))
	want := []int{0, 10, 20}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got := ToSlice(MapWhile(Empty[int](), func(x int) (int, bool) { return x, true })); len(got) != 0 {
		t.Fatalf("got %v, want empty", got)
	}
}

func TestFlatMap(t *testing.T) {
	got := ToSlice(FlatMap(Range1(3), func(x int) iter.Seq[int] {
		return Range1(x + 1)
	}))
	want := []int{0, 0, 1, 0, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestFlatten(t *testing.T) {
	src := seqOf[iter.Seq[int]](Range1(2), Range2(10, 12))
	got := ToSlice(Flatten(src))
	want := []int{0, 1, 10, 11}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSplit(t *testing.T) {
	src := seqOf(1, 2)
	got := ToMap(Split[int, int](src, func(e int) (int, int) { return e, e * 10 }))
	want := map[int]int{1: 10, 2: 20}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCast(t *testing.T) {
	src := seqOf[any](int(1), "not int", int(3))
	var vals []int
	var oks []bool
	for v, ok := range Cast[int](src) {
		vals = append(vals, v)
		oks = append(oks, ok)
	}
	wantVals := []int{1, 0, 3}
	wantOks := []bool{true, false, true}
	if !reflect.DeepEqual(vals, wantVals) {
		t.Fatalf("vals got %v, want %v", vals, wantVals)
	}
	if !reflect.DeepEqual(oks, wantOks) {
		t.Fatalf("oks got %v, want %v", oks, wantOks)
	}
}

func TestChain(t *testing.T) {
	got := ToSlice(Chain(Range1(2), Range2(10, 12)))
	want := []int{0, 1, 10, 11}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got := ToSlice(Chain(Empty[int](), Empty[int]())); len(got) != 0 {
		t.Fatalf("got %v, want empty", got)
	}
}

func TestForEach(t *testing.T) {
	sum := 0
	ForEach(Range1(5), func(v int) { sum += v })
	if sum != 10 {
		t.Fatalf("got %d, want 10", sum)
	}
}

func TestTryForEachNoError(t *testing.T) {
	visited := 0
	err := TryForEach(Range1(3), func(v int) error { visited++; return nil })
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if visited != 3 {
		t.Fatalf("visited %d, want 3", visited)
	}
}

func TestTryReduceNoError(t *testing.T) {
	got, ok, err := TryReduce(Range1(5), func(acc, v int) (int, error) { return acc + v, nil })
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if !ok {
		t.Fatal("ok=false, want true")
	}
	if got != 10 {
		t.Fatalf("got %d, want 10", got)
	}
}

func TestTryFoldNoError(t *testing.T) {
	got, err := TryFold(Range1(5), 0, func(acc, v int) (int, error) { return acc + v, nil })
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if got != 10 {
		t.Fatalf("got %d, want 10", got)
	}
}

func TestSizeValue(t *testing.T) {
	src := seqOf(1, 2, 1)
	if n := SizeValue(src, 1); n != 2 {
		t.Fatalf("got %d, want 2", n)
	}
}

func TestContains(t *testing.T) {
	if !Contains(Range1(5), 3) {
		t.Fatal("Contains should be true for 3")
	}
	if Contains(Range1(5), 10) {
		t.Fatal("Contains should be false for 10")
	}
	if Contains(Empty[int](), 0) {
		t.Fatal("Contains should be false for empty")
	}
}

func TestContainsFunc(t *testing.T) {
	if !ContainsFunc(Range1(5), func(x int) bool { return x == 3 }) {
		t.Fatal("ContainsFunc should be true")
	}
	if ContainsFunc(Range1(5), func(x int) bool { return x == 99 }) {
		t.Fatal("ContainsFunc should be false")
	}
}

func TestFirst(t *testing.T) {
	v, ok := First(Range1(5))
	if !ok || v != 0 {
		t.Fatalf("got (%d, %v), want (0, true)", v, ok)
	}
	_, ok = First(Empty[int]())
	if ok {
		t.Fatal("First on empty should return false")
	}
}

func TestFirstFunc(t *testing.T) {
	v, ok := FirstFunc(Range1(5), func(x int) bool { return x > 2 })
	if !ok || v != 3 {
		t.Fatalf("got (%d, %v), want (3, true)", v, ok)
	}
	_, ok = FirstFunc(Range1(5), func(x int) bool { return x > 100 })
	if ok {
		t.Fatal("FirstFunc no match should return false")
	}
}

func TestLast(t *testing.T) {
	v, ok := Last(Range1(5))
	if !ok || v != 4 {
		t.Fatalf("got (%d, %v), want (4, true)", v, ok)
	}
	_, ok = Last(Empty[int]())
	if ok {
		t.Fatal("Last on empty should return false")
	}
}

func TestLastFunc(t *testing.T) {
	src := seqOf(1, 2, 3)
	v, ok := LastFunc(src, func(x int) bool { return x < 3 })
	if !ok || v != 2 {
		t.Fatalf("got (%d, %v), want (2, true)", v, ok)
	}
	_, ok = LastFunc(src, func(x int) bool { return x > 100 })
	if ok {
		t.Fatal("LastFunc no match should return false")
	}
}

func TestPosition(t *testing.T) {
	idx, ok := Position(Range1(5), func(x int) bool { return x == 3 })
	if !ok || idx != 3 {
		t.Fatalf("got (%d, %v), want (3, true)", idx, ok)
	}
	_, ok = Position(Range1(5), func(x int) bool { return x == 100 })
	if ok {
		t.Fatal("Position no match should return false")
	}
}

func TestNth(t *testing.T) {
	cases := []struct {
		name string
		src  iter.Seq[int]
		n    int
		want int
		ok   bool
	}{
		{"normal", Range1(10), 3, 3, true},
		{"first", Range1(10), 0, 0, true},
		{"last", Range1(5), 4, 4, true},
		{"out_of_range", Range1(2), 5, 0, false},
		{"negative", Range1(5), -1, 0, false},
		{"empty", Empty[int](), 0, 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, ok := Nth(c.src, c.n)
			if ok != c.ok || (ok && v != c.want) {
				t.Fatalf("got (%d, %v), want (%d, %v)", v, ok, c.want, c.ok)
			}
		})
	}
}

func TestFindMap(t *testing.T) {
	parse := func(s string) (int, bool) {
		var n int
		var neg bool
		i := 0
		if i < len(s) && s[i] == '-' {
			neg = true
			i++
		}
		if i >= len(s) {
			return 0, false
		}
		for ; i < len(s); i++ {
			if s[i] < '0' || s[i] > '9' {
				return 0, false
			}
			n = n*10 + int(s[i]-'0')
		}
		if neg {
			n = -n
		}
		return n, true
	}

	v, ok := FindMap(seqOf("x", "12", "y"), parse)
	if !ok || v != 12 {
		t.Fatalf("got (%d, %v), want (12, true)", v, ok)
	}

	_, ok = FindMap(seqOf("x", "y"), parse)
	if ok {
		t.Fatal("FindMap no match should return false")
	}

	_, ok = FindMap(Empty[string](), parse)
	if ok {
		t.Fatal("FindMap on empty should return false")
	}
}

func TestEqual(t *testing.T) {
	if !Equal(Range1(3), Range1(3)) {
		t.Fatal("equal sequences should be true")
	}
	if Equal(Range1(3), Range1(4)) {
		t.Fatal("different length should be false")
	}
	if Equal(Range1(3), Range2(1, 4)) {
		t.Fatal("different values should be false")
	}
	if Equal(Range1(3), Empty[int]()) {
		t.Fatal("non-empty vs empty should be false")
	}
}

func TestEqualFunc(t *testing.T) {
	cmp := func(a, b int) bool { return a == b }
	if !EqualFunc(Range1(3), Range1(3), cmp) {
		t.Fatal("equal sequences should be true")
	}
	if EqualFunc(Range1(3), Range1(4), cmp) {
		t.Fatal("different length should be false")
	}
}

func TestMax(t *testing.T) {
	v, ok := Max(seqOf(3, 1, 5, 4))
	if !ok || v != 5 {
		t.Fatalf("got (%d, %v), want (5, true)", v, ok)
	}
	_, ok = Max(Empty[int]())
	if ok {
		t.Fatal("Max on empty should return false")
	}
}

func TestMaxFunc(t *testing.T) {
	v, ok := MaxFunc(seqOf(3, 1, 5), func(a, b int) int { return a - b })
	if !ok || v != 5 {
		t.Fatalf("got (%d, %v), want (5, true)", v, ok)
	}
	_, ok = MaxFunc(Empty[int](), func(a, b int) int { return a - b })
	if ok {
		t.Fatal("MaxFunc on empty should return false")
	}
}

func TestMin(t *testing.T) {
	v, ok := Min(seqOf(3, -1, 4))
	if !ok || v != -1 {
		t.Fatalf("got (%d, %v), want (-1, true)", v, ok)
	}
	_, ok = Min(Empty[int]())
	if ok {
		t.Fatal("Min on empty should return false")
	}
}

func TestMinFunc(t *testing.T) {
	v, ok := MinFunc(seqOf(3, -1, 4), func(a, b int) int { return a - b })
	if !ok || v != -1 {
		t.Fatalf("got (%d, %v), want (-1, true)", v, ok)
	}
	_, ok = MinFunc(Empty[int](), func(a, b int) int { return a - b })
	if ok {
		t.Fatal("MinFunc on empty should return false")
	}
}

func TestIsSorted(t *testing.T) {
	if !IsSorted(Range1(5)) {
		t.Fatal("ascending should be sorted")
	}
	if !IsSorted(seqOf(5, 4, 3)) {
		t.Fatal("descending consistent order should be sorted")
	}
	if IsSorted(seqOf(1, 3, 2)) {
		t.Fatal("unsorted should be false")
	}
	if !IsSorted(Empty[int]()) {
		t.Fatal("empty should be sorted")
	}
	if !IsSorted(Once(5)) {
		t.Fatal("single should be sorted")
	}
}

func TestIsSortedFunc(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	if !IsSortedFunc(Range1(5), cmp) {
		t.Fatal("ascending should be sorted")
	}
	if IsSortedFunc(seqOf(1, 3, 2), cmp) {
		t.Fatal("unsorted should be false")
	}
	if !IsSortedFunc(Empty[int](), cmp) {
		t.Fatal("empty should be sorted")
	}
}

// TestGeneratorsEarlyStop triggers the !yield early-return branch in all
// generator-style functions.
func TestGeneratorsEarlyStop(t *testing.T) {
	stopEarly(Range1(10))
	stopEarly(Range2(1, 10))
	stopEarly(Range3(1, 10, 2))
	stopEarly(Range3(10, 1, -2))
	stopEarly(FromFunc(func() (int, bool) { return 1, true }))
	stopEarly(Once(1))
	stopEarly(Repeat(1))
	stopEarly(Map(Range1(10), func(x int) int { return x }))
	stopEarly(MapWhile(Range1(10), func(x int) (int, bool) { return x, true }))
	stopEarly(FlatMap(Range1(10), func(x int) iter.Seq[int] { return Once(x) }))
	stopEarly(Flatten(func(yield func(iter.Seq[int]) bool) { yield(Range1(10)) }))
	stopEarly(Inspect(Range1(10), func(int) {}))
	stopEarly(Scan(Range1(10), 0, func(a, e int) (int, bool) { return a + e, true }))
	stopEarly(Filter(Range1(10), func(int) bool { return true }))
	stopEarly(FilterMap(Range1(10), func(x int) (int, bool) { return x, true }))
	stopEarly(Take(Range1(10), 5))
	stopEarly(TakeWhile(Range1(10), func(int) bool { return true }))
	stopEarly(Skip(Range1(10), 2))
	stopEarly(SkipWhile(Range1(10), func(int) bool { return false }))
	stopEarly(StepBy(Range1(10), 2))
	stopEarly(Chain(Range1(10), Range1(10)))
	stopEarly(Chain(Empty[int](), Range1(10))) // second segment !yield branch
	stopEarly(ZipWith(Range1(10), Range1(10), func(a, b int) int { return a }))
	// Seq2-returning generators in seq.go
	stopEarly2(Enumerate(Range1(10)))
	stopEarly2(Split[int, int](Range1(10), func(e int) (int, int) { return e, e }))
	stopEarly2(Cast[int](func(yield func(any) bool) { yield(1) }))
	stopEarly2(Zip(Range1(10), Range1(10)))
}
