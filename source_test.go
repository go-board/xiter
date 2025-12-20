package xiter

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange1(t *testing.T) {
	// Test the normal case.
	seq := Range1(5)
	result := ToSlice(seq)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)

	// Test boundary case.
	seq = Range1(0)
	result = ToSlice(seq)
	assert.Empty(t, result)
}

func TestRange2(t *testing.T) {
	// Test the normal case.
	seq := Range2(2, 6)
	result := ToSlice(seq)
	assert.Equal(t, []int{2, 3, 4, 5}, result)

	// Test boundary case.
	seq = Range2(5, 3)
	result = ToSlice(seq)
	assert.Empty(t, result)
}

func TestRange3(t *testing.T) {
	// Test the normal case.
	seq := Range3(1, 10, 2)
	result := ToSlice(seq)
	assert.Equal(t, []int{1, 3, 5, 7, 9}, result)

	// Test boundary case.
	seq = Range3(5, 5, 1)
	result = ToSlice(seq)
	assert.Empty(t, result)
}

func TestFromFunc(t *testing.T) {
	// Test the normal case.
	count := 0
	seq := FromFunc(func() (int, bool) {
		count++
		return count, count <= 3
	})
	result := ToSlice(seq)
	assert.Equal(t, []int{1, 2, 3}, result)

	// Test immediate termination.
	seq = FromFunc(func() (int, bool) {
		return 0, false
	})
	result = ToSlice(seq)
	assert.Empty(t, result)
}

func TestFromFunc2(t *testing.T) {
	// Test the normal case.
	count := 0
	seq := FromFunc2(func() (int, string, bool) {
		count++
		return count, "val" + strconv.Itoa(count), count <= 2
	})
	result := ToSlice2(seq)
	assert.Equal(t, []Pair[int, string]{{1, "val1"}, {2, "val2"}}, result)
}

func TestFromSlice(t *testing.T) {
	// Test the normal case.
	slice := []string{"a", "b", "c"}
	seq := fromSlice(slice)
	result := ToSlice(seq)
	assert.Equal(t, []string{"a", "b", "c"}, result)

	// Test empty slice.
	seq2 := fromSlice([]int{})
	result2 := ToSlice(seq2)
	assert.Empty(t, result2)
}

func TestFromSlice2(t *testing.T) {
	// Test the normal case.
	slice := []Pair[int, string]{{1, "x"}, {2, "y"}}
	seq := FromSlice2(slice)
	result := ToSlice2(seq)
	assert.Equal(t, []Pair[int, string]{{1, "x"}, {2, "y"}}, result)
}

func TestFromMap(t *testing.T) {
	// Test the normal case.
	m := map[string]int{"a": 1, "b": 2}
	seq := FromMap(m)
	result := ToSlice2(seq)
	// Map iteration order is non-deterministic, so only verify existence.
	assert.Len(t, result, 2)
	vals := make(map[string]int)
	for _, p := range result {
		vals[p.Key] = p.Value
	}
	assert.Equal(t, 1, vals["a"])
	assert.Equal(t, 2, vals["b"])

	// Test empty map.
	seq2 := FromMap(map[int]int{})
	result2 := ToSlice2(seq2)
	assert.Empty(t, result2)
}

func TestOnce(t *testing.T) {
	// Test the normal case.
	seq := Once("hello")
	result := ToSlice(seq)
	assert.Equal(t, []string{"hello"}, result)
}

func TestOnce2(t *testing.T) {
	// Test the normal case.
	seq := Once2(1, "x")
	result := ToSlice2(seq)
	assert.Equal(t, []Pair[int, string]{{1, "x"}}, result)
}

func TestEmpty(t *testing.T) {
	// Test empty sequence.
	seq := Empty[int]()
	result := ToSlice(seq)
	assert.Empty(t, result)
}

func TestEmpty2(t *testing.T) {
	// Test empty key/value sequence.
	seq := Empty2[int, string]()
	result := ToSlice2(seq)
	assert.Empty(t, result)
}

func TestRepeat(t *testing.T) {
	// Test repeating sequence (paired with Take).
	seq := Repeat("a")
	limited := Take(seq, 3)
	result := ToSlice(limited)
	assert.Equal(t, []string{"a", "a", "a"}, result)
}

func TestRepeat2(t *testing.T) {
	// Test repeating key/value sequence (paired with Take2).
	seq := Repeat2(1, "x")
	limited := Take2(seq, 2)
	result := ToSlice2(limited)
	assert.Equal(t, []Pair[int, string]{{1, "x"}, {1, "x"}}, result)
}
