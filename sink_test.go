package xiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSlice(t *testing.T) {
	// 测试整数序列转换
	seq := Map(Range1(4), func(n int) int { return n + 1 })
	result := ToSlice(seq)
	assert.Equal(t, []int{1, 2, 3, 4}, result)

	// 测试字符串序列转换
	strSeq := fromSlice([]string{"a", "b", "c"})
	strResult := ToSlice(strSeq)
	assert.Equal(t, []string{"a", "b", "c"}, strResult)

	// 测试空序列
	emptySeq := Empty[int]()
	emptyResult := ToSlice(emptySeq)
	assert.Empty(t, emptyResult)
}

func TestToMap(t *testing.T) {
	// 测试正常情况
	pairs := FromSlice2([]Pair[string, int]{{"a", 1}, {"b", 2}})
	result := ToMap(pairs)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, result)

	// 测试重复键（后面的会覆盖前面的）
	dupPairs := FromSlice2([]Pair[string, int]{{"a", 1}, {"a", 10}, {"b", 2}})
	dupResult := ToMap(dupPairs)
	assert.Equal(t, map[string]int{"a": 10, "b": 2}, dupResult)

	// 测试空序列
	emptySeq := Empty2[string, int]()
	emptyResult := ToMap(emptySeq)
	assert.Empty(t, emptyResult)
}

func TestToSet(t *testing.T) {
	// 测试包含重复元素的序列
	seq := fromSlice([]int{1, 2, 2, 3, 3, 3})
	result := ToSet(seq)
	assert.Equal(t, map[int]struct{}{1: {}, 2: {}, 3: {}}, result)

	// 测试空序列
	emptySeq := Empty[string]()
	emptyResult := ToSet(emptySeq)
	assert.Empty(t, emptyResult)
}

func TestToSlice2(t *testing.T) {
	// 测试正常情况
	seq := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	result := ToSlice2(seq)
	assert.Equal(t, []Pair[int, string]{{1, "a"}, {2, "b"}}, result)

	// 测试空序列
	emptySeq := Empty2[int, string]()
	emptyResult := ToSlice2(emptySeq)
	assert.Empty(t, emptyResult)
}

func TestGroupBy(t *testing.T) {
	// 测试按首字母分组字符串
	words := fromSlice([]string{"apple", "banana", "apricot", "cherry", "avocado"})
	groups := GroupBy(words, func(s string) rune { return []rune(s)[0] })
	assert.Equal(t, []string{"apple", "apricot", "avocado"}, groups['a'])
	assert.Equal(t, []string{"banana"}, groups['b'])
	assert.Equal(t, []string{"cherry"}, groups['c'])

	// 测试按奇偶分组整数
	numbers := Range1(5)
	numGroups := GroupBy(numbers, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	assert.Equal(t, []int{0, 2, 4}, numGroups["even"])
	assert.Equal(t, []int{1, 3}, numGroups["odd"])

	// 测试空序列
	emptySeq := Empty[int]()
	emptyGroups := GroupBy(emptySeq, func(n int) int { return n })
	assert.Empty(t, emptyGroups)
}
