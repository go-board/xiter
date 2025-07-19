package xiter

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange1(t *testing.T) {
	// 测试正常情况
	seq := Range1(5)
	result := ToSlice(seq)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)

	// 测试边界情况
	seq = Range1(0)
	result = ToSlice(seq)
	assert.Empty(t, result)
}

func TestRange2(t *testing.T) {
	// 测试正常情况
	seq := Range2(2, 6)
	result := ToSlice(seq)
	assert.Equal(t, []int{2, 3, 4, 5}, result)

	// 测试边界情况
	seq = Range2(5, 3)
	result = ToSlice(seq)
	assert.Empty(t, result)
}

func TestRange3(t *testing.T) {
	// 测试正常情况
	seq := Range3(1, 10, 2)
	result := ToSlice(seq)
	assert.Equal(t, []int{1, 3, 5, 7, 9}, result)

	// 测试边界情况
	seq = Range3(5, 5, 1)
	result = ToSlice(seq)
	assert.Empty(t, result)
}

func TestFromFunc(t *testing.T) {
	// 测试正常情况
	count := 0
	seq := FromFunc(func() (int, bool) {
		count++
		return count, count <= 3
	})
	result := ToSlice(seq)
	assert.Equal(t, []int{1, 2, 3}, result)

	// 测试立即结束
	seq = FromFunc(func() (int, bool) {
		return 0, false
	})
	result = ToSlice(seq)
	assert.Empty(t, result)
}

func TestFromFunc2(t *testing.T) {
	// 测试正常情况
	count := 0
	seq := FromFunc2(func() (int, string, bool) {
		count++
		return count, "val" + strconv.Itoa(count), count <= 2
	})
	result := ToSlice2(seq)
	assert.Equal(t, []Pair[int, string]{{1, "val1"}, {2, "val2"}}, result)
}

func TestFromSlice(t *testing.T) {
	// 测试正常情况
	slice := []string{"a", "b", "c"}
	seq := fromSlice(slice)
	result := ToSlice(seq)
	assert.Equal(t, []string{"a", "b", "c"}, result)

	// 测试空切片
	seq2 := fromSlice([]int{})
	result2 := ToSlice(seq2)
	assert.Empty(t, result2)
}

func TestFromSlice2(t *testing.T) {
	// 测试正常情况
	slice := []Pair[int, string]{{1, "x"}, {2, "y"}}
	seq := FromSlice2(slice)
	result := ToSlice2(seq)
	assert.Equal(t, []Pair[int, string]{{1, "x"}, {2, "y"}}, result)
}

func TestFromMap(t *testing.T) {
	// 测试正常情况
	m := map[string]int{"a": 1, "b": 2}
	seq := FromMap(m)
	result := ToSlice2(seq)
	// 由于map遍历顺序不确定，我们只检查元素是否存在
	assert.Len(t, result, 2)
	vals := make(map[string]int)
	for _, p := range result {
		vals[p.Key] = p.Value
	}
	assert.Equal(t, 1, vals["a"])
	assert.Equal(t, 2, vals["b"])

	// 测试空map
	seq2 := FromMap(map[int]int{})
	result2 := ToSlice2(seq2)
	assert.Empty(t, result2)
}

func TestOnce(t *testing.T) {
	// 测试正常情况
	seq := Once("hello")
	result := ToSlice(seq)
	assert.Equal(t, []string{"hello"}, result)
}

func TestOnce2(t *testing.T) {
	// 测试正常情况
	seq := Once2(1, "x")
	result := ToSlice2(seq)
	assert.Equal(t, []Pair[int, string]{{1, "x"}}, result)
}

func TestEmpty(t *testing.T) {
	// 测试空序列
	seq := Empty[int]()
	result := ToSlice(seq)
	assert.Empty(t, result)
}

func TestEmpty2(t *testing.T) {
	// 测试空键值对序列
	seq := Empty2[int, string]()
	result := ToSlice2(seq)
	assert.Empty(t, result)
}

func TestRepeat(t *testing.T) {
	// 测试重复序列（配合Take使用）
	seq := Repeat("a")
	limited := Take(seq, 3)
	result := ToSlice(limited)
	assert.Equal(t, []string{"a", "a", "a"}, result)
}

func TestRepeat2(t *testing.T) {
	// 测试重复键值对序列（配合Take2使用）
	seq := Repeat2(1, "x")
	limited := Take2(seq, 2)
	result := ToSlice2(limited)
	assert.Equal(t, []Pair[int, string]{{1, "x"}, {1, "x"}}, result)
}
