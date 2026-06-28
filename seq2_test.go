package xiter

import (
	"errors"
	"reflect"
	"slices"
	"testing"
)

// ============================================================================
// Source
// ============================================================================

func TestFromFunc2(t *testing.T) {
	count := 0
	s := FromFunc2(func() (string, int, bool) {
		count++
		if count <= 2 {
			return "key" + string(rune('0'+count)), count * 10, true
		}
		return "", 0, false
	})
	result := ToMap(s)
	expected := map[string]int{"key1": 10, "key2": 20}
	if len(result) != len(expected) {
		t.Errorf("FromFunc2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("FromFunc2[%s] = %d, expected %d", k, result[k], v)
		}
	}

	s = FromFunc2(func() (string, int, bool) {
		return "", 0, false
	})
	result = ToMap(s)
	if len(result) != 0 {
		t.Errorf("FromFunc2 returned %d elements, expected 0", len(result))
	}
}

func TestIterate2(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		s := Iterate2(0, 1, func(k, v int) (int, int, bool) {
			if k >= 3 {
				return 0, 0, false
			}
			return k + 1, v * 2, true
		})
		var gotK []int
		var gotV []int
		s(func(k, v int) bool {
			gotK = append(gotK, k)
			gotV = append(gotV, v)
			return true
		})
		if !slices.Equal(gotK, []int{0, 1, 2, 3}) || !slices.Equal(gotV, []int{1, 2, 4, 8}) {
			t.Fatalf("got keys=%v values=%v", gotK, gotV)
		}
	})
	t.Run("seed only when next immediately stops", func(t *testing.T) {
		s := Iterate2("a", 1, func(k string, v int) (string, int, bool) { return "", 0, false })
		var gotK []string
		var gotV []int
		s(func(k string, v int) bool {
			gotK = append(gotK, k)
			gotV = append(gotV, v)
			return true
		})
		if len(gotK) != 1 || gotK[0] != "a" || gotV[0] != 1 {
			t.Fatalf("got keys=%v values=%v", gotK, gotV)
		}
	})
	t.Run("early stop", func(t *testing.T) {
		// Break after 2 yields; Iterate2 must honor the stop signal.
		var nextCalls int
		s := Iterate2(0, 0, func(k, v int) (int, int, bool) {
			nextCalls++
			if k >= 1 {
				return 0, 0, false
			}
			return k + 1, v + 1, true
		})
		var gotK []int
		var gotV []int
		n := 0
		s(func(k, v int) bool {
			n++
			if n > 2 {
				return false
			}
			gotK = append(gotK, k)
			gotV = append(gotV, v)
			return true
		})
		if len(gotK) != 2 {
			t.Fatalf("got %d pairs, want 2", len(gotK))
		}
		if nextCalls == 0 {
			t.Fatalf("next was never called")
		}
	})
	t.Run("yield false on seed", func(t *testing.T) {
		// Consumer rejects the very first pair (the seed): next must never be
		// called.
		var nextCalls int
		s := Iterate2(0, 0, func(k, v int) (int, int, bool) {
			nextCalls++
			return k + 1, v + 1, true
		})
		s(func(k, v int) bool { return false })
		if nextCalls != 0 {
			t.Fatalf("next was called %d times, want 0", nextCalls)
		}
	})
	t.Run("yield false after first pair", func(t *testing.T) {
		// Consumer accepts the seed but rejects the second pair: next must not
		// be called again.
		var nextCalls int
		s := Iterate2(0, 0, func(k, v int) (int, int, bool) {
			nextCalls++
			return k + 1, v + 1, true
		})
		n := 0
		s(func(k, v int) bool {
			n++
			return n == 1
		})
		if nextCalls != 1 {
			t.Fatalf("next was called %d times, want 1", nextCalls)
		}
	})
}

func TestOnce2(t *testing.T) {
	s := Once2("key", 42)
	result := ToMap(s)
	if len(result) != 1 {
		t.Errorf("Once2 returned %d elements, expected 1", len(result))
	}
	if result["key"] != 42 {
		t.Errorf("Once2 returned %d, expected 42", result["key"])
	}
}

func TestEmpty2(t *testing.T) {
	s := Empty2[string, int]()
	result := ToMap(s)
	if len(result) != 0 {
		t.Errorf("Empty2 returned %d elements, expected 0", len(result))
	}
}

func TestRepeat2(t *testing.T) {
	s := Repeat2("key", "value")
	taken := Take2(s, 3)
	result := ToMap(taken)
	// 由于是重复的键值对，最终map中只会有一个元素
	if len(result) != 1 {
		t.Errorf("Repeat2 returned %d elements, expected 1", len(result))
	}
	if result["key"] != "value" {
		t.Errorf("Repeat2 returned %s, expected %s", result["key"], "value")
	}
}

// ============================================================================
// Transform
// ============================================================================

func TestMap2(t *testing.T) {
	s := Enumerate(Range1(5))
	mapped := Map2(s, func(k, v int) (int, int) { return k, v * 2 })
	result := ToMap(mapped)
	expected := map[int]int{0: 0, 1: 2, 2: 4, 3: 6, 4: 8}
	if len(result) != len(expected) {
		t.Errorf("Map2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Map2[%d] = %d, expected %d", k, result[k], v)
		}
	}
}

func TestInspect2(t *testing.T) {
	seen := map[int]int{}
	got := ToMap(Inspect2(Enumerate(Range1(3)), func(k, v int) {
		seen[k] = v
	}))
	want := map[int]int{0: 0, 1: 1, 2: 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if !reflect.DeepEqual(seen, want) {
		t.Fatalf("seen %v, want %v", seen, want)
	}
}

func TestKeys(t *testing.T) {
	got := ToSlice(Keys(Enumerate(Range2(10, 13))))
	want := []int{0, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestValues(t *testing.T) {
	got := ToSlice(Values(Enumerate(Range2(10, 13))))
	want := []int{10, 11, 12}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSwap(t *testing.T) {
	got := ToMap(Swap(Enumerate(Range2(10, 13))))
	want := map[int]int{10: 0, 11: 1, 12: 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// ============================================================================
// Filter / Slice
// ============================================================================

func TestTake2(t *testing.T) {
	s := Enumerate(Range1(10))
	taken := Take2(s, 5)
	result := ToMap(taken)
	expected := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4}
	if len(result) != len(expected) {
		t.Errorf("Take2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Take2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	s = Enumerate(Range1(5))
	taken = Take2(s, 10)
	result = ToMap(taken)
	expected = map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4}
	if len(result) != len(expected) {
		t.Errorf("Take2 returned %d elements, expected %d", len(result), len(expected))
	}

	emptySeq := Empty2[int, int]()
	taken = Take2(emptySeq, 5)
	result = ToMap(taken)
	if len(result) != 0 {
		t.Errorf("Take2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestSkip2(t *testing.T) {
	s := Enumerate(Range1(10))
	skiped := Skip2(s, 5)
	result := ToMap(skiped)
	expected := map[int]int{5: 5, 6: 6, 7: 7, 8: 8, 9: 9}
	if len(result) != len(expected) {
		t.Errorf("Skip2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Skip2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	s = Enumerate(Range1(5))
	skiped = Skip2(s, 10)
	result = ToMap(skiped)
	if len(result) != 0 {
		t.Errorf("Skip2 returned %d elements, expected 0", len(result))
	}

	emptySeq := Empty2[int, int]()
	skiped = Skip2(emptySeq, 5)
	result = ToMap(skiped)
	if len(result) != 0 {
		t.Errorf("Skip2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestFilter2(t *testing.T) {
	s := Enumerate(Range1(10))
	filtered := Filter2(s, func(k, v int) bool { return k%2 == 0 })
	result := ToMap(filtered)
	expected := map[int]int{0: 0, 2: 2, 4: 4, 6: 6, 8: 8}
	if len(result) != len(expected) {
		t.Errorf("Filter2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Filter2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	emptySeq := Empty2[int, int]()
	filtered = Filter2(emptySeq, func(k, v int) bool { return true })
	result = ToMap(filtered)
	if len(result) != 0 {
		t.Errorf("Filter2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestFilterMap2(t *testing.T) {
	source := func(yield func(int, int) bool) {
		for i := 1; i <= 5; i++ {
			if !yield(i, i*10) {
				return
			}
		}
	}
	filterFn := func(k, v int) (int, int, bool) {
		if k%2 == 0 {
			return k, v * 2, true
		}
		return 0, 0, false
	}
	result := FilterMap2(source, filterFn)
	expected := map[int]int{2: 40, 4: 80}
	actual := ToMap(result)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FilterMap2 failed. Expected: %v, Got: %v", expected, actual)
	}
}

func TestTakeWhile2(t *testing.T) {
	s := Enumerate(Range1(10))
	takenWhile := TakeWhile2(s, func(k, v int) bool { return k <= 5 })
	result := ToMap(takenWhile)
	expected := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	if len(result) != len(expected) {
		t.Errorf("TakeWhile2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("TakeWhile2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	takenWhile = TakeWhile2(s, func(k, v int) bool { return false })
	result = ToMap(takenWhile)
	if len(result) != 0 {
		t.Errorf("TakeWhile2 returned %d elements, expected 0", len(result))
	}

	emptySeq := Empty2[int, int]()
	takenWhile = TakeWhile2(emptySeq, func(k, v int) bool { return true })
	result = ToMap(takenWhile)
	if len(result) != 0 {
		t.Errorf("TakeWhile2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestSkipWhile2(t *testing.T) {
	s := Enumerate(Range1(10))
	skipedWhile := SkipWhile2(s, func(k, v int) bool { return k <= 5 })
	result := ToMap(skipedWhile)
	expected := map[int]int{6: 6, 7: 7, 8: 8, 9: 9}
	if len(result) != len(expected) {
		t.Errorf("SkipWhile2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("SkipWhile2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	skipedWhile = SkipWhile2(s, func(k, v int) bool { return false })
	result = ToMap(skipedWhile)
	expected = map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9}
	if len(result) != len(expected) {
		t.Errorf("SkipWhile2 returned %d elements, expected %d", len(result), len(expected))
	}

	emptySeq := Empty2[int, int]()
	skipedWhile = SkipWhile2(emptySeq, func(k, v int) bool { return true })
	result = ToMap(skipedWhile)
	if len(result) != 0 {
		t.Errorf("SkipWhile2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

// ============================================================================
// Terminal
// ============================================================================

func TestFold2(t *testing.T) {
	s := Enumerate(Range1(5))
	sum := Fold2(s, 0, func(acc, k, v int) int { return acc + k + v })
	expectedSum := 20 // (0+0)+(1+1)+(2+2)+(3+3)+(4+4)
	if sum != expectedSum {
		t.Errorf("Fold2 sum: expected %d, got %d", expectedSum, sum)
	}

	emptySeq := Empty2[int, int]()
	emptyResult := Fold2(emptySeq, 42, func(acc, k, v int) int { return acc + k + v })
	if emptyResult != 42 {
		t.Errorf("Fold2 on empty sequence: expected 42, got %d", emptyResult)
	}
}

func TestReduce2(t *testing.T) {
	s := Enumerate(Range1(5))
	key, val, ok := Reduce2(s, func(k1, v1, k2, v2 int) (int, int) {
		return k1 + k2, v1 + v2
	})
	if !ok {
		t.Fatal("Reduce2 returned ok=false, want true")
	}
	if key != 10 || val != 10 {
		t.Fatalf("got (%d, %d), want (10, 10)", key, val)
	}
}

func TestReduce2Empty(t *testing.T) {
	s := Empty2[int, int]()
	key, val, ok := Reduce2(s, func(k1, v1, k2, v2 int) (int, int) {
		return k1 + k2, v1 + v2
	})
	if ok {
		t.Fatal("Reduce2 returned ok=true, want false")
	}
	if key != 0 || val != 0 {
		t.Fatalf("got (%d, %d), want zero values", key, val)
	}
}

func TestTryReduce2(t *testing.T) {
	wantErr := errors.New("stop")
	s := Enumerate(Range1(5))
	key, val, ok, err := TryReduce2(s, func(k1, v1, k2, v2 int) (int, int, error) {
		if k2 == 3 {
			return k1 + k2, v1 + v2, wantErr
		}
		return k1 + k2, v1 + v2, nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if !ok {
		t.Fatal("TryReduce2 returned ok=false, want true")
	}
	// 0+1+2+3=6
	if key != 6 || val != 6 {
		t.Fatalf("got (%d, %d), want (6, 6)", key, val)
	}
}

func TestTryForEach2(t *testing.T) {
	wantErr := errors.New("stop")
	visited := 0
	err := TryForEach2(Enumerate(Range1(5)), func(k, v int) error {
		visited++
		if k == 2 || v == 2 {
			return wantErr
		}
		return nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if visited != 3 {
		t.Fatalf("visited %d pairs, want 3", visited)
	}
}

func TestTryFold2(t *testing.T) {
	wantErr := errors.New("stop")
	got, err := TryFold2(Enumerate(Range1(5)), 0, func(acc, k, v int) (int, error) {
		acc += k + v
		if k == 3 {
			return acc, wantErr
		}
		return acc, nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if got != 12 {
		t.Fatalf("got accumulator %d, want 12", got)
	}
}

func TestSize2(t *testing.T) {
	s := Enumerate(Range1(5))
	size := Size2(s)
	expected := 5
	if size != expected {
		t.Errorf("Size2: expected %d, got %d", expected, size)
	}

	emptySeq := Empty2[int, int]()
	emptySize := Size2(emptySeq)
	if emptySize != 0 {
		t.Errorf("Size2 on empty sequence: expected 0, got %d", emptySize)
	}

	infinite := Repeat2("key", 1)
	limited := Take2(infinite, 10)
	limitedSize := Size2(limited)
	if limitedSize != 10 {
		t.Errorf("Size2 on limited infinite sequence: expected 10, got %d", limitedSize)
	}
}

func TestSizeFunc2(t *testing.T) {
	s := Enumerate(Range1(10))
	evenKeyCount := SizeFunc2(s, func(k, v int) bool { return k%2 == 0 })
	expected := 5
	if evenKeyCount != expected {
		t.Errorf("SizeFunc2 (even keys): expected %d, got %d", expected, evenKeyCount)
	}

	emptySeq := Empty2[int, int]()
	emptyCount := SizeFunc2(emptySeq, func(k, v int) bool { return true })
	if emptyCount != 0 {
		t.Errorf("SizeFunc2 on empty sequence: expected 0, got %d", emptyCount)
	}
}

func TestSizeValue2(t *testing.T) {
	s := func(yield func(int, int) bool) {
		yield(1, 10)
		yield(2, 20)
		yield(1, 10)
		yield(3, 30)
	}
	count := SizeValue2(s, 1, 10)
	if count != 2 {
		t.Errorf("SizeValue2: expected 2, got %d", count)
	}
}

// ============================================================================
// Compare / Search
// ============================================================================

func TestAny2(t *testing.T) {
	s := Enumerate(Range1(10))
	result := Any2(s, func(k, v int) bool { return k == 5 })
	if !result {
		t.Errorf("Any2 should return true when key == 5 exists")
	}

	s = Enumerate(Range1(10))
	result = Any2(s, func(k, v int) bool { return k > 100 })
	if result {
		t.Errorf("Any2 should return false for keys > 100")
	}

	empty := Empty2[int, int]()
	result = Any2(empty, func(k, v int) bool { return true })
	if result {
		t.Errorf("Any2 should return false for empty sequence")
	}
}

func TestAll2(t *testing.T) {
	s := Enumerate(Range1(10))
	result := All2(s, func(k, v int) bool { return k == v })
	if !result {
		t.Errorf("All2 should return true when all keys equal values")
	}

	s = Enumerate(Range1(10))
	result = All2(s, func(k, v int) bool { return k < 5 })
	if result {
		t.Errorf("All2 should return false when not all keys < 5")
	}

	empty := Empty2[int, int]()
	result = All2(empty, func(k, v int) bool { return true })
	if !result {
		t.Errorf("All2 should return true for empty sequence")
	}
}

func TestCompare2(t *testing.T) {
	s1 := Enumerate(Range1(3))
	s2 := Enumerate(Range1(3))
	result := Compare2(s1, s2)
	if result != 0 {
		t.Errorf("Compare2 (equal sequences): expected 0, got %d", result)
	}

	s3 := Enumerate(Range1(2))
	s4 := Enumerate(Range1(3))
	result = Compare2(s3, s4)
	if result >= 0 {
		t.Errorf("Compare2 (s1 < s2): expected < 0, got %d", result)
	}

	s5 := Enumerate(Range1(3))
	s6 := Enumerate(Range1(2))
	result = Compare2(s5, s6)
	if result <= 0 {
		t.Errorf("Compare2 (s1 > s2): expected > 0, got %d", result)
	}
}

func TestCompareFunc2(t *testing.T) {
	s1 := func(yield func(int, int) bool) {
		yield(1, 10)
		yield(2, 20)
	}
	s2 := func(yield func(int, int) bool) {
		yield(1, 10)
		yield(2, 20)
	}
	result := CompareFunc2(s1, s2, func(k1, v1, k2, v2 int) int {
		if k1 != k2 {
			return k1 - k2
		}
		return v1 - v2
	})
	if result != 0 {
		t.Errorf("CompareFunc2 (equal sequences): expected 0, got %d", result)
	}

	s3 := func(yield func(int, int) bool) {
		yield(1, 10)
		yield(2, 20)
	}
	s4 := func(yield func(int, int) bool) {
		yield(1, 30)
		yield(2, 20)
	}
	result = CompareFunc2(s3, s4, func(k1, v1, k2, v2 int) int {
		if k1 != k2 {
			return k1 - k2
		}
		return v1 - v2
	})
	if result >= 0 {
		t.Errorf("CompareFunc2 (v1 < v2): expected < 0, got %d", result)
	}
}

// ============================================================================
// Supplementary coverage
// ============================================================================

func TestMapWhile2(t *testing.T) {
	got := ToMap(MapWhile2(Enumerate(Range1(5)), func(k, v int) (int, int, bool) {
		if k < 2 {
			return k, v * 10, true
		}
		return 0, 0, false
	}))
	want := map[int]int{0: 0, 1: 10}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestChain2(t *testing.T) {
	src1 := func(yield func(int, int) bool) {
		if !yield(1, 1) {
			return
		}
		if !yield(2, 2) {
			return
		}
	}
	src2 := func(yield func(int, int) bool) {
		if !yield(3, 3) {
			return
		}
	}
	got := ToMap(Chain2(src1, src2))
	want := map[int]int{1: 1, 2: 2, 3: 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got := ToMap(Chain2(Empty2[int, int](), Empty2[int, int]())); len(got) != 0 {
		t.Fatalf("got %v, want empty", got)
	}
}

func TestForEach2(t *testing.T) {
	sum := 0
	ForEach2(Enumerate(Range1(3)), func(k, v int) { sum += k + v })
	if sum != 6 {
		t.Fatalf("got %d, want 6", sum)
	}
}

func TestTryForEach2NoError(t *testing.T) {
	visited := 0
	err := TryForEach2(Enumerate(Range1(3)), func(k, v int) error { visited++; return nil })
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if visited != 3 {
		t.Fatalf("visited %d, want 3", visited)
	}
}

func TestTryReduce2NoError(t *testing.T) {
	k, v, ok, err := TryReduce2(Enumerate(Range1(5)), func(k1, v1, k2, v2 int) (int, int, error) {
		return k1 + k2, v1 + v2, nil
	})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if !ok {
		t.Fatal("ok=false, want true")
	}
	if k != 10 || v != 10 {
		t.Fatalf("got (%d, %d), want (10, 10)", k, v)
	}
}

func TestTryReduce2Empty(t *testing.T) {
	k, v, ok, err := TryReduce2(Empty2[int, int](), func(k1, v1, k2, v2 int) (int, int, error) {
		return k1 + k2, v1 + v2, nil
	})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if ok {
		t.Fatal("ok=true, want false")
	}
	if k != 0 || v != 0 {
		t.Fatalf("got (%d, %d), want (0, 0)", k, v)
	}
}

func TestTryFold2NoError(t *testing.T) {
	got, err := TryFold2(Enumerate(Range1(5)), 0, func(acc, k, v int) (int, error) { return acc + k + v, nil })
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if got != 20 {
		t.Fatalf("got %d, want 20", got)
	}
}

func TestContains2(t *testing.T) {
	src := func(yield func(int, int) bool) {
		if !yield(1, 10) {
			return
		}
		if !yield(2, 20) {
			return
		}
	}
	if !Contains2(src, 1, 10) {
		t.Fatal("Contains2 should be true for (1,10)")
	}
	if Contains2(src, 1, 20) {
		t.Fatal("Contains2 should be false for (1,20)")
	}
	if Contains2(Empty2[int, int](), 0, 0) {
		t.Fatal("Contains2 should be false for empty")
	}
}

func TestContainsFunc2(t *testing.T) {
	src := func(yield func(int, int) bool) {
		if !yield(1, 10) {
			return
		}
		if !yield(2, 20) {
			return
		}
	}
	if !ContainsFunc2(src, func(k, v int) bool { return v == 20 }) {
		t.Fatal("ContainsFunc2 should be true for v==20")
	}
	if ContainsFunc2(src, func(k, v int) bool { return v == 999 }) {
		t.Fatal("ContainsFunc2 should be false for v==999")
	}
}

func TestFirst2(t *testing.T) {
	k, v, ok := First2(Enumerate(Range1(5)))
	if !ok || k != 0 || v != 0 {
		t.Fatalf("got (%d, %d, %v), want (0, 0, true)", k, v, ok)
	}
	_, _, ok = First2(Empty2[int, int]())
	if ok {
		t.Fatal("First2 on empty should return false")
	}
}

func TestFirstFunc2(t *testing.T) {
	k, v, ok := FirstFunc2(Enumerate(Range1(5)), func(k, v int) bool { return k == 2 })
	if !ok || k != 2 || v != 2 {
		t.Fatalf("got (%d, %d, %v), want (2, 2, true)", k, v, ok)
	}
	_, _, ok = FirstFunc2(Enumerate(Range1(5)), func(k, v int) bool { return k == 99 })
	if ok {
		t.Fatal("FirstFunc2 no match should return false")
	}
}

func TestLast2(t *testing.T) {
	k, v, ok := Last2(Enumerate(Range1(5)))
	if !ok || k != 4 || v != 4 {
		t.Fatalf("got (%d, %d, %v), want (4, 4, true)", k, v, ok)
	}
	_, _, ok = Last2(Empty2[int, int]())
	if ok {
		t.Fatal("Last2 on empty should return false")
	}
}

func TestLastFunc2(t *testing.T) {
	src := func(yield func(int, int) bool) {
		if !yield(0, 0) {
			return
		}
		if !yield(1, 1) {
			return
		}
		if !yield(2, 2) {
			return
		}
	}
	k, v, ok := LastFunc2(src, func(k, v int) bool { return k < 2 })
	if !ok || k != 1 || v != 1 {
		t.Fatalf("got (%d, %d, %v), want (1, 1, true)", k, v, ok)
	}
	_, _, ok = LastFunc2(src, func(k, v int) bool { return k > 99 })
	if ok {
		t.Fatal("LastFunc2 no match should return false")
	}
}

func TestPosition2(t *testing.T) {
	idx, ok := Position2(Enumerate(Range1(5)), func(k, v int) bool { return k == 3 })
	if !ok || idx != 3 {
		t.Fatalf("got (%d, %v), want (3, true)", idx, ok)
	}
	_, ok = Position2(Enumerate(Range1(5)), func(k, v int) bool { return k == 99 })
	if ok {
		t.Fatal("Position2 no match should return false")
	}
}

func TestEqual2(t *testing.T) {
	if !Equal2(Enumerate(Range1(3)), Enumerate(Range1(3))) {
		t.Fatal("equal sequences should be true")
	}
	if Equal2(Enumerate(Range1(3)), Enumerate(Range1(4))) {
		t.Fatal("different length should be false")
	}
	if Equal2(Enumerate(Range1(3)), Empty2[int, int]()) {
		t.Fatal("non-empty vs empty should be false")
	}
}

func TestEqualFunc2(t *testing.T) {
	cmp := func(k1, v1, k2, v2 int) bool { return k1 == k2 && v1 == v2 }
	if !EqualFunc2(Enumerate(Range1(3)), Enumerate(Range1(3)), cmp) {
		t.Fatal("equal sequences should be true")
	}
	if EqualFunc2(Enumerate(Range1(3)), Enumerate(Range1(4)), cmp) {
		t.Fatal("different length should be false")
	}
}

func TestCompare2KeyDiff(t *testing.T) {
	x := Once2(1, 10)
	yGreater := Once2(2, 10)
	if c := Compare2(x, yGreater); c >= 0 {
		t.Fatalf("x<y should be <0, got %d", c)
	}
	yLesser := Once2(0, 10)
	if c := Compare2(x, yLesser); c <= 0 {
		t.Fatalf("x>y should be >0, got %d", c)
	}
	// value differs when key equal
	a := Once2(1, 10)
	b := Once2(1, 20)
	if c := Compare2(a, b); c >= 0 {
		t.Fatalf("a<b by value should be <0, got %d", c)
	}
}

// TestGeneratorsEarlyStop2 triggers the !yield early-return branch in all
// Seq2 generator-style functions.
func TestGeneratorsEarlyStop2(t *testing.T) {
	stopEarly2(FromFunc2(func() (int, int, bool) { return 1, 2, true }))
	stopEarly2(Map2(Enumerate(Range1(10)), func(k, v int) (int, int) { return k, v }))
	stopEarly2(MapWhile2(Enumerate(Range1(10)), func(k, v int) (int, int, bool) { return k, v, true }))
	stopEarly2(Inspect2(Enumerate(Range1(10)), func(int, int) {}))
	stopEarly2(Swap(Enumerate(Range1(10))))
	stopEarly2(Filter2(Enumerate(Range1(10)), func(int, int) bool { return true }))
	stopEarly2(FilterMap2(Enumerate(Range1(10)), func(k, v int) (int, int, bool) { return k, v, true }))
	stopEarly2(Take2(Enumerate(Range1(10)), 5))
	stopEarly2(Take2(Enumerate(Range1(10)), 0)) // n<=0 early return
	stopEarly2(TakeWhile2(Enumerate(Range1(10)), func(int, int) bool { return true }))
	stopEarly2(Skip2(Enumerate(Range1(10)), 2))
	stopEarly2(SkipWhile2(Enumerate(Range1(10)), func(int, int) bool { return false }))
	stopEarly2(Chain2(Enumerate(Range1(10)), Enumerate(Range1(10))))
	stopEarly2(Chain2(Empty2[int, int](), Enumerate(Range1(10)))) // second segment !yield branch
	// Seq-returning generators in seq2.go
	stopEarly(Join(Enumerate(Range1(10)), func(k, v int) int { return v }))
	stopEarly(Keys(Enumerate(Range1(10))))
	stopEarly(Values(Enumerate(Range1(10))))
}
