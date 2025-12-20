package xiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSlice(t *testing.T) {
	// Test integer sequence conversion.
	seq := Map(Range1(4), func(n int) int { return n + 1 })
	result := ToSlice(seq)
	assert.Equal(t, []int{1, 2, 3, 4}, result)

	// Test string sequence conversion.
	strSeq := fromSlice([]string{"a", "b", "c"})
	strResult := ToSlice(strSeq)
	assert.Equal(t, []string{"a", "b", "c"}, strResult)

	// Test empty sequence.
	emptySeq := Empty[int]()
	emptyResult := ToSlice(emptySeq)
	assert.Empty(t, emptyResult)
}

func TestToMap(t *testing.T) {
	// Test the normal case.
	pairs := FromSlice2([]Pair[string, int]{{"a", 1}, {"b", 2}})
	result := ToMap(pairs)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, result)

	// Test duplicate keys (later entries overwrite earlier ones).
	dupPairs := FromSlice2([]Pair[string, int]{{"a", 1}, {"a", 10}, {"b", 2}})
	dupResult := ToMap(dupPairs)
	assert.Equal(t, map[string]int{"a": 10, "b": 2}, dupResult)

	// Test empty sequence.
	emptySeq := Empty2[string, int]()
	emptyResult := ToMap(emptySeq)
	assert.Empty(t, emptyResult)
}

func TestToSet(t *testing.T) {
	// Test sequence with duplicate elements.
	seq := fromSlice([]int{1, 2, 2, 3, 3, 3})
	result := ToSet(seq)
	assert.Equal(t, map[int]struct{}{1: {}, 2: {}, 3: {}}, result)

	// Test empty sequence.
	emptySeq := Empty[string]()
	emptyResult := ToSet(emptySeq)
	assert.Empty(t, emptyResult)
}

func TestToSlice2(t *testing.T) {
	// Test the normal case.
	seq := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	result := ToSlice2(seq)
	assert.Equal(t, []Pair[int, string]{{1, "a"}, {2, "b"}}, result)

	// Test empty sequence.
	emptySeq := Empty2[int, string]()
	emptyResult := ToSlice2(emptySeq)
	assert.Empty(t, emptyResult)
}

func TestGroupBy(t *testing.T) {
	// Test grouping strings by first letter.
	words := fromSlice([]string{"apple", "banana", "apricot", "cherry", "avocado"})
	groups := GroupBy(words, func(s string) rune { return []rune(s)[0] })
	assert.Equal(t, []string{"apple", "apricot", "avocado"}, groups['a'])
	assert.Equal(t, []string{"banana"}, groups['b'])
	assert.Equal(t, []string{"cherry"}, groups['c'])

	// Test grouping integers by parity.
	numbers := Range1(5)
	numGroups := GroupBy(numbers, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	assert.Equal(t, []int{0, 2, 4}, numGroups["even"])
	assert.Equal(t, []int{1, 3}, numGroups["odd"])

	// Test empty sequence.
	emptySeq := Empty[int]()
	emptyGroups := GroupBy(emptySeq, func(n int) int { return n })
	assert.Empty(t, emptyGroups)
}
