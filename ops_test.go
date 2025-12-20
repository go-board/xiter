package xiter

import (
	"cmp"
	"fmt"
	"iter"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fromSlice creates a helper that returns an iter.Seq.
func fromSlice[E any](s []E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func TestEnumerate(t *testing.T) {
	// Test the normal case.
	seq := Range1(3)
	enumerated := Enumerate(seq)
	result := ToSlice2(enumerated)
	assert.Equal(t, []Pair[int, int]{{0, 0}, {1, 1}, {2, 2}}, result)

	// Test empty sequence.
	seq = Empty[int]()
	enumerated = Enumerate(seq)
	result = ToSlice2(enumerated)
	assert.Empty(t, result)
}

func TestMap(t *testing.T) {
	// Test doubling integers.
	input := Range1(4)
	doubled := Map(input, func(e int) int { return e * 2 })
	result := ToSlice(doubled)
	assert.Equal(t, []int{0, 2, 4, 6}, result)

	// Test type conversion.
	numbers := Range1(3)
	strings := Map(numbers, func(n int) string { return fmt.Sprintf("%c", 'A'+n) })
	resultStr := ToSlice(strings)
	assert.Equal(t, []string{"A", "B", "C"}, resultStr)
}

func TestFilter(t *testing.T) {
	// Test filtering even numbers.
	input := Range1(5)
	evens := Filter(input, func(e int) bool { return e%2 == 0 })
	result := ToSlice(evens)
	assert.Equal(t, []int{0, 2, 4}, result)

	// Test no matches.
	odds := Filter(input, func(e int) bool { return e < 0 })
	result = ToSlice(odds)
	assert.Empty(t, result)
}

func TestSkip(t *testing.T) {
	// Test skipping elements.
	input := Range1(5)
	skipped := Skip(input, 2)
	result := ToSlice(skipped)
	assert.Equal(t, []int{2, 3, 4}, result)
}

func TestSkip2(t *testing.T) {
	// Test skipping key/value pairs.
	input := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}})
	skipped := Skip2(input, 1)
	result := ToSlice2(skipped)
	assert.Equal(t, []Pair[int, string]{{2, "b"}, {3, "c"}}, result)
}

func TestTakeWhile(t *testing.T) {
	// Test taking while predicate holds.
	input := Range1(5)
	taken := TakeWhile(input, func(e int) bool { return e < 3 })
	result := ToSlice(taken)
	assert.Equal(t, []int{0, 1, 2}, result)
}

func TestTakeWhile2(t *testing.T) {
	// Test taking key/value pairs while predicate holds.
	input := FromSlice2([]Pair[int, int]{{1, 1}, {2, 2}, {3, 3}})
	taken := TakeWhile2(input, func(k, v int) bool { return v < 3 })
	result := ToSlice2(taken)
	assert.Equal(t, []Pair[int, int]{{1, 1}, {2, 2}}, result)
}

func TestSkipWhile(t *testing.T) {
	// Test skipping while predicate holds.
	input := Range1(5)
	skipped := SkipWhile(input, func(e int) bool { return e < 3 })
	result := ToSlice(skipped)
	assert.Equal(t, []int{3, 4}, result)
}

func TestSkipWhile2(t *testing.T) {
	// Test skipping key/value pairs while predicate holds.
	input := FromSlice2([]Pair[int, int]{{1, 1}, {2, 2}, {3, 3}})
	skipped := SkipWhile2(input, func(k, v int) bool { return v < 3 })
	result := ToSlice2(skipped)
	assert.Equal(t, []Pair[int, int]{{3, 3}}, result)
}

func TestFind(t *testing.T) {
	// Test finding the first matching element.
	input := Range1(5)
	found, ok := Find(input, func(e int) bool { return e == 3 })
	assert.True(t, ok)
	assert.Equal(t, 3, found)
}

func TestFind2(t *testing.T) {
	// Test finding the first matching key/value pair.
	input := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}})
	k, v, ok := Find2(input, func(k int, v string) bool { return v == "b" })
	assert.True(t, ok)
	assert.Equal(t, 2, k)
	assert.Equal(t, "b", v)
}

func TestPosition(t *testing.T) {
	// Test locating the index of a matching element.
	input := Range1(5)
	index, ok := Position(input, func(e int) bool { return e == 4 })
	assert.True(t, ok)
	assert.Equal(t, 4, index)
}

func TestPosition2(t *testing.T) {
	// Test locating the index of a matching key/value pair.
	input := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}})
	index, ok := Position2(input, func(k int, v string) bool { return v == "c" })
	assert.True(t, ok)
	assert.Equal(t, 2, index)
}

func TestFold(t *testing.T) {
	// Test summation.
	numbers := Range1(5)
	sum := Fold(numbers, 0, func(acc, e int) int { return acc + e })
	assert.Equal(t, 10, sum)

	// Test string concatenation.
	words := fromSlice([]string{"Hello", " ", "World"})
	phrase := Fold(words, "", func(acc, s string) string { return acc + s })
	assert.Equal(t, "Hello World", phrase)
}

func TestSize(t *testing.T) {
	// Test non-empty sequence.
	seq := Range1(5)
	assert.Equal(t, 5, Size(seq))

	// Test empty sequence.
	seq = Empty[int]()
	assert.Equal(t, 0, Size(seq))
}

func TestMax(t *testing.T) {
	// Test the normal case.
	numbers := fromSlice([]int{3, 1, 4, 1, 5, 9})
	maxVal, ok := Max(numbers)
	assert.True(t, ok)
	assert.Equal(t, 9, maxVal)

	// Test empty sequence.
	emptySeq := Empty[int]()
	maxVal, ok = Max(emptySeq)
	assert.False(t, ok)
	assert.Equal(t, 0, maxVal)
}

func TestMin(t *testing.T) {
	// Test the normal case.
	numbers := fromSlice([]int{5, 2, 7, 1, 3})
	minVal, ok := Min(numbers)
	assert.True(t, ok)
	assert.Equal(t, 1, minVal)
}

func TestEqual(t *testing.T) {
	// Test equal sequences.
	seq1 := fromSlice([]int{1, 2, 3})
	seq2 := fromSlice([]int{1, 2, 3})
	assert.True(t, Equal(seq1, seq2))

	// Test unequal sequences.
	seq3 := fromSlice([]int{1, 3, 2})
	assert.False(t, Equal(seq1, seq3))
}

func TestSum(t *testing.T) {
	// Test summing integers.
	numbers := Range1(5)
	assert.Equal(t, 10, Sum(numbers))

	// Test empty sequence.
	emptySeq := Empty[int]()
	assert.Equal(t, 0, Sum(emptySeq))
}

func TestAny(t *testing.T) {
	// Test when a matching element exists.
	numbers := Range1(5)
	assert.True(t, Any(numbers, func(e int) bool { return e == 3 }))

	// Test when no matching element exists.
	assert.False(t, Any(numbers, func(e int) bool { return e > 10 }))
}

func TestAll(t *testing.T) {
	// Test when all elements match.
	numbers := fromSlice([]int{2, 4, 6, 8})
	assert.True(t, All(numbers, func(e int) bool { return e > 0 }))

	// Test when a non-matching element exists.
	assert.False(t, All(numbers, func(e int) bool { return e < 5 }))
}

func TestConcat(t *testing.T) {
	// Test concatenating multiple sequences.
	seq1 := Range1(3)
	seq2 := Range2(5, 7)
	combined := Concat(seq1, seq2)
	result := ToSlice(combined)
	assert.Equal(t, []int{0, 1, 2, 5, 6}, result)
}

func TestConcat2(t *testing.T) {
	// Test concatenating multiple key/value sequences.
	seq1 := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	seq2 := FromSlice2([]Pair[int, string]{{3, "c"}})
	combined := Concat2(seq1, seq2)
	result := ToSlice2(combined)
	assert.Equal(t, []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}, result)

	// Test concatenating empty sequences.
	emptySeq := Empty2[int, string]()
	combined = Concat2(emptySeq, emptySeq)
	result = ToSlice2(combined)
	assert.Empty(t, result)
}

func TestKeys(t *testing.T) {
	// Test extracting keys.
	pairs := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}})
	keys := Keys(pairs)
	result := ToSlice(keys)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestValues(t *testing.T) {
	// Test extracting values.
	pairs := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}})
	values := Values(pairs)
	result := ToSlice(values)
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestEqualFunc(t *testing.T) {
	// Test custom comparison function.
	seq1 := fromSlice([]string{"Apple", "Banana"})
	seq2 := fromSlice([]string{"apple", "banana"})
	assert.True(t, EqualFunc(seq1, seq2, func(a, b string) bool {
		return strings.EqualFold(a, b)
	}))

	// Test unequal case.
	seq3 := fromSlice([]string{"apple"})
	assert.False(t, EqualFunc(seq1, seq3, func(a, b string) bool {
		return strings.EqualFold(a, b)
	}))
}

func TestEqual2(t *testing.T) {
	// Test equal key/value sequences.
	seq1 := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	seq2 := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	assert.True(t, Equal2(seq1, seq2))

	// Test unequal key/value sequences.
	seq3 := FromSlice2([]Pair[int, string]{{1, "a"}, {3, "b"}})
	assert.False(t, Equal2(seq1, seq3))
}

func TestJoin(t *testing.T) {
	// Test joining key/value pairs into strings.
	pairs := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	strs := Join(pairs, func(k int, v string) string { return fmt.Sprintf("%d=%s", k, v) })
	result := ToSlice(strs)
	assert.Equal(t, []string{"1=a", "2=b"}, result)
}

func TestSplit(t *testing.T) {
	// Test splitting strings into key/value pairs.
	strs := fromSlice([]string{"1=a", "2=b"})
	pairs := Split(strs, func(s string) (int, string) {
		parts := strings.Split(s, "=")
		k, _ := strconv.Atoi(parts[0])
		return k, parts[1]
	})
	result := ToSlice2(pairs)
	assert.Equal(t, []Pair[int, string]{{1, "a"}, {2, "b"}}, result)
}

func TestDistinct(t *testing.T) {
	// Test de-duplication.
	seq := fromSlice([]int{1, 2, 2, 3, 3, 3})
	distinct := Distinct(seq)
	result := ToSlice(distinct)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestFilterMap(t *testing.T) {
	// Test filtering and mapping.
	numbers := Range1(5)
	squares := FilterMap(numbers, func(e int) (int, bool) {
		if e%2 == 0 {
			return e * e, true
		}
		return 0, false
	})
	result := ToSlice(squares)
	assert.Equal(t, []int{0, 4, 16}, result)
}

func TestFilterMap2(t *testing.T) {
	// Test filtering and mapping key/value pairs.
	input := FromSlice2([]Pair[int, int]{{1, 10}, {2, 20}, {3, 30}})
	filtered := FilterMap2(input, func(k, v int) (string, int, bool) {
		if v > 15 {
			return fmt.Sprintf("key%d", k), v / 10, true
		}
		return "", 0, false
	})
	result := ToSlice2(filtered)
	assert.Equal(t, []Pair[string, int]{{"key2", 2}, {"key3", 3}}, result)
}

func TestForEach(t *testing.T) {
	// Test side-effect operations.
	var sum int
	seq := Range1(5)
	ForEach(seq, func(e int) { sum += e })
	assert.Equal(t, 10, sum)
}

func TestForEach2(t *testing.T) {
	// Test side-effect operations on key/value pairs.
	var keys []int
	input := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	ForEach2(input, func(k int, v string) { keys = append(keys, k) })
	assert.Equal(t, []int{1, 2}, keys)
}

func TestSizeFunc(t *testing.T) {
	// Test conditional counting.
	numbers := Range1(10)
	count := SizeFunc(numbers, func(e int) bool { return e%3 == 0 })
	assert.Equal(t, 4, count) // 0,3,6,9
}

func TestSizeFunc2(t *testing.T) {
	// Test conditional counting for key/value pairs.
	input := FromSlice2([]Pair[int, int]{{1, 5}, {2, 10}, {3, 15}})
	count := SizeFunc2(input, func(k, v int) bool { return v > 10 })
	assert.Equal(t, 1, count)
}

func TestMaxFunc(t *testing.T) {
	// Test custom comparison function.
	words := fromSlice([]string{"apple", "banana", "cherry"})
	longest, ok := MaxFunc(words, func(a, b string) int {
		return len(a) - len(b)
	})
	assert.True(t, ok)
	assert.Equal(t, "banana", longest)
}

func TestMinFunc(t *testing.T) {
	// Test custom comparison function.
	words := fromSlice([]string{"apple", "banana", "cherry", "date"})
	shortest, ok := MinFunc(words, func(a, b string) int {
		return len(a) - len(b)
	})
	assert.True(t, ok)
	assert.Equal(t, "date", shortest)
}

func TestContainsFunc(t *testing.T) {
	// Test custom containment predicate.
	words := fromSlice([]string{"apple", "banana", "cherry"})
	assert.True(t, ContainsFunc(words, func(s string) bool { return len(s) > 5 }))
	assert.False(t, ContainsFunc(words, func(s string) bool { return len(s) > 10 }))
}

func TestContains2(t *testing.T) {
	// Test containment for key/value pairs.
	input := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	assert.True(t, Contains2(input, 2, "b"))
	assert.False(t, Contains2(input, 3, "c"))
}

func TestIsSortedFunc(t *testing.T) {
	// Test custom sorted check.
	numbers := fromSlice([]int{1, 3, 5, 7})
	assert.True(t, IsSortedFunc(numbers, cmp.Compare[int]))
	shuffled := fromSlice([]int{3, 1, 5, 7})
	assert.False(t, IsSortedFunc(shuffled, cmp.Compare[int]))
}

func TestCast(t *testing.T) {
	// Test type conversion.
	values := fromSlice([]any{1, "two", 3, "four"})
	casted := Cast[int](values)
	result := ToSlice2(casted)
	expected := []Pair[int, bool]{{1, true}, {0, false}, {3, true}, {0, false}}
	assert.Equal(t, expected, result)
}

func TestSizeValue(t *testing.T) {
	// Test counting values.
	numbers := fromSlice([]int{2, 5, 5, 7, 5})
	assert.Equal(t, 3, SizeValue(numbers, 5))
	assert.Equal(t, 0, SizeValue(numbers, 10))
}

func TestSizeValue2(t *testing.T) {
	// Test counting values in key/value pairs.
	pairs := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "a"}})
	assert.Equal(t, 2, SizeValue2(pairs, "a"))
}

func TestFold2(t *testing.T) {
	// Test folding key/value pairs.
	input := FromSlice2([]Pair[int, int]{{1, 10}, {2, 20}, {3, 30}})
	sum := Fold2(input, 0, func(acc, k, v int) int { return acc + v })
	assert.Equal(t, 60, sum)
}

func TestAny2(t *testing.T) {
	// Test existence in key/value pairs.
	input := FromSlice2([]Pair[int, int]{{1, 5}, {2, 15}, {3, 25}})
	assert.True(t, Any2(input, func(k, v int) bool { return v > 20 }))
	assert.False(t, Any2(input, func(k, v int) bool { return v > 30 }))
}

func TestAll2(t *testing.T) {
	// Test all-matching check for key/value pairs.
	input := FromSlice2([]Pair[int, int]{{1, 5}, {2, 15}, {3, 25}})
	assert.True(t, All2(input, func(k, v int) bool { return v > 0 }))
	assert.False(t, All2(input, func(k, v int) bool { return v < 20 }))
}

func TestEqualFunc2(t *testing.T) {
	// Test custom comparison for key/value pairs.
	seq1 := FromSlice2([]Pair[string, string]{{"a", "Apple"}, {"b", "Banana"}})
	seq2 := FromSlice2([]Pair[string, string]{{"a", "apple"}, {"b", "banana"}})
	assert.True(t, EqualFunc2(seq1, seq2, func(k1, v1, k2, v2 string) bool {
		return k1 == k2 && strings.EqualFold(v1, v2)
	}))
}

func TestMap2(t *testing.T) {
	// Test mapping key/value pairs.
	input := FromSlice2([]Pair[int, int]{{1, 10}, {2, 20}, {3, 30}})
	mapped := Map2(input, func(k, v int) (string, int) { return fmt.Sprintf("key%d", k), v * 2 })
	result := ToSlice2(mapped)
	assert.Equal(t, []Pair[string, int]{{"key1", 20}, {"key2", 40}, {"key3", 60}}, result)
}

func TestFilter2(t *testing.T) {
	// Test filtering key/value pairs.
	input := FromSlice2([]Pair[int, int]{{1, 5}, {2, 15}, {3, 25}, {4, 5}})
	filtered := Filter2(input, func(k, v int) bool { return v > 10 })
	result := ToSlice2(filtered)
	assert.Equal(t, []Pair[int, int]{{2, 15}, {3, 25}}, result)
}

func TestContains(t *testing.T) {
	// Test contains element.
	seq := fromSlice([]string{"a", "b", "c"})
	assert.True(t, Contains(seq, "b"))

	// Test does not contain element.
	assert.False(t, Contains(seq, "d"))
}

func TestIsSorted(t *testing.T) {
	// Test sorted sequence.
	seq := fromSlice([]int{1, 2, 3, 4})
	assert.True(t, IsSorted(seq))

	// Test unsorted sequence.
	seq = fromSlice([]int{1, 3, 2, 4})
	assert.False(t, IsSorted(seq))
}
