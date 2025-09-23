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

// fromSlice 创建一个返回iter.Seq的辅助函数
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
	// 测试正常情况
	seq := Range1(3)
	enumerated := Enumerate(seq)
	result := ToSlice2(enumerated)
	assert.Equal(t, []Pair[int, int]{{0, 0}, {1, 1}, {2, 2}}, result)

	// 测试空序列
	seq = Empty[int]()
	enumerated = Enumerate(seq)
	result = ToSlice2(enumerated)
	assert.Empty(t, result)
}

func TestMap(t *testing.T) {
	// 测试整数加倍
	input := Range1(4)
	doubled := Map(input, func(e int) int { return e * 2 })
	result := ToSlice(doubled)
	assert.Equal(t, []int{0, 2, 4, 6}, result)

	// 测试类型转换
	numbers := Range1(3)
	strings := Map(numbers, func(n int) string { return fmt.Sprintf("%c", 'A'+n) })
	resultStr := ToSlice(strings)
	assert.Equal(t, []string{"A", "B", "C"}, resultStr)
}

func TestFilter(t *testing.T) {
	// 测试筛选偶数
	input := Range1(5)
	evens := Filter(input, func(e int) bool { return e%2 == 0 })
	result := ToSlice(evens)
	assert.Equal(t, []int{0, 2, 4}, result)

	// 测试无匹配项
	odds := Filter(input, func(e int) bool { return e < 0 })
	result = ToSlice(odds)
	assert.Empty(t, result)
}

func TestFold(t *testing.T) {
	// 测试求和
	numbers := Range1(5)
	sum := Fold(numbers, 0, func(acc, e int) int { return acc + e })
	assert.Equal(t, 10, sum)

	// 测试字符串拼接
	words := fromSlice([]string{"Hello", " ", "World"})
	phrase := Fold(words, "", func(acc, s string) string { return acc + s })
	assert.Equal(t, "Hello World", phrase)
}

func TestSize(t *testing.T) {
	// 测试非空序列
	seq := Range1(5)
	assert.Equal(t, 5, Size(seq))

	// 测试空序列
	seq = Empty[int]()
	assert.Equal(t, 0, Size(seq))
}

func TestMax(t *testing.T) {
	// 测试正常情况
	numbers := fromSlice([]int{3, 1, 4, 1, 5, 9})
	maxVal, ok := Max(numbers)
	assert.True(t, ok)
	assert.Equal(t, 9, maxVal)

	// 测试空序列
	emptySeq := Empty[int]()
	maxVal, ok = Max(emptySeq)
	assert.False(t, ok)
	assert.Equal(t, 0, maxVal)
}

func TestMin(t *testing.T) {
	// 测试正常情况
	numbers := fromSlice([]int{5, 2, 7, 1, 3})
	minVal, ok := Min(numbers)
	assert.True(t, ok)
	assert.Equal(t, 1, minVal)
}

func TestEqual(t *testing.T) {
	// 测试相等序列
	seq1 := fromSlice([]int{1, 2, 3})
	seq2 := fromSlice([]int{1, 2, 3})
	assert.True(t, Equal(seq1, seq2))

	// 测试不相等序列
	seq3 := fromSlice([]int{1, 3, 2})
	assert.False(t, Equal(seq1, seq3))
}

func TestSum(t *testing.T) {
	// 测试整数求和
	numbers := Range1(5)
	assert.Equal(t, 10, Sum(numbers))

	// 测试空序列
	emptySeq := Empty[int]()
	assert.Equal(t, 0, Sum(emptySeq))
}

func TestAny(t *testing.T) {
	// 测试存在匹配元素
	numbers := Range1(5)
	assert.True(t, Any(numbers, func(e int) bool { return e == 3 }))

	// 测试不存在匹配元素
	assert.False(t, Any(numbers, func(e int) bool { return e > 10 }))
}

func TestAll(t *testing.T) {
	// 测试所有元素匹配
	numbers := fromSlice([]int{2, 4, 6, 8})
	assert.True(t, All(numbers, func(e int) bool { return e > 0 }))

	// 测试存在不匹配元素
	assert.False(t, All(numbers, func(e int) bool { return e < 5 }))
}

func TestConcat(t *testing.T) {
	// 测试连接多个序列
	seq1 := Range1(3)
	seq2 := Range2(5, 7)
	combined := Concat(seq1, seq2)
	result := ToSlice(combined)
	assert.Equal(t, []int{0, 1, 2, 5, 6}, result)
}

func TestConcat2(t *testing.T) {
	// 测试连接多个键值对序列
	seq1 := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	seq2 := FromSlice2([]Pair[int, string]{{3, "c"}})
	combined := Concat2(seq1, seq2)
	result := ToSlice2(combined)
	assert.Equal(t, []Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}, result)

	// 测试空序列连接
	emptySeq := Empty2[int, string]()
	combined = Concat2(emptySeq, emptySeq)
	result = ToSlice2(combined)
	assert.Empty(t, result)
}

func TestKeys(t *testing.T) {
	// 测试提取键
	pairs := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}})
	keys := Keys(pairs)
	result := ToSlice(keys)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestValues(t *testing.T) {
	// 测试提取值
	pairs := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}})
	values := Values(pairs)
	result := ToSlice(values)
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestEqualFunc(t *testing.T) {
	// 测试自定义比较函数
	seq1 := fromSlice([]string{"Apple", "Banana"})
	seq2 := fromSlice([]string{"apple", "banana"})
	assert.True(t, EqualFunc(seq1, seq2, func(a, b string) bool {
		return strings.EqualFold(a, b)
	}))

	// 测试不相等情况
	seq3 := fromSlice([]string{"apple"})
	assert.False(t, EqualFunc(seq1, seq3, func(a, b string) bool {
		return strings.EqualFold(a, b)
	}))
}

func TestEqual2(t *testing.T) {
	// 测试键值对序列相等
	seq1 := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	seq2 := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	assert.True(t, Equal2(seq1, seq2))

	// 测试键值对序列不相等
	seq3 := FromSlice2([]Pair[int, string]{{1, "a"}, {3, "b"}})
	assert.False(t, Equal2(seq1, seq3))
}

func TestJoin(t *testing.T) {
	// 测试连接键值对为字符串
	pairs := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	strs := Join(pairs, func(k int, v string) string { return fmt.Sprintf("%d=%s", k, v) })
	result := ToSlice(strs)
	assert.Equal(t, []string{"1=a", "2=b"}, result)
}

func TestSplit(t *testing.T) {
	// 测试拆分字符串为键值对
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
	// 测试去重
	seq := fromSlice([]int{1, 2, 2, 3, 3, 3})
	distinct := Distinct(seq)
	result := ToSlice(distinct)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestFilterMap(t *testing.T) {
	// 测试筛选并转换
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
	// 测试键值对筛选并转换
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
	// 测试副作用操作
	var sum int
	seq := Range1(5)
	ForEach(seq, func(e int) { sum += e })
	assert.Equal(t, 10, sum)
}

func TestForEach2(t *testing.T) {
	// 测试键值对副作用操作
	var keys []int
	input := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	ForEach2(input, func(k int, v string) { keys = append(keys, k) })
	assert.Equal(t, []int{1, 2}, keys)
}

func TestSizeFunc(t *testing.T) {
	// 测试条件计数
	numbers := Range1(10)
	count := SizeFunc(numbers, func(e int) bool { return e%3 == 0 })
	assert.Equal(t, 4, count) // 0,3,6,9
}

func TestSizeFunc2(t *testing.T) {
	// 测试键值对条件计数
	input := FromSlice2([]Pair[int, int]{{1, 5}, {2, 10}, {3, 15}})
	count := SizeFunc2(input, func(k, v int) bool { return v > 10 })
	assert.Equal(t, 1, count)
}

func TestMaxFunc(t *testing.T) {
	// 测试自定义比较函数
	words := fromSlice([]string{"apple", "banana", "cherry"})
	longest, ok := MaxFunc(words, func(a, b string) int {
		return len(a) - len(b)
	})
	assert.True(t, ok)
	assert.Equal(t, "banana", longest)
}

func TestMinFunc(t *testing.T) {
	// 测试自定义比较函数
	words := fromSlice([]string{"apple", "banana", "cherry", "date"})
	shortest, ok := MinFunc(words, func(a, b string) int {
		return len(a) - len(b)
	})
	assert.True(t, ok)
	assert.Equal(t, "date", shortest)
}

func TestContainsFunc(t *testing.T) {
	// 测试自定义包含判断
	words := fromSlice([]string{"apple", "banana", "cherry"})
	assert.True(t, ContainsFunc(words, func(s string) bool { return len(s) > 5 }))
	assert.False(t, ContainsFunc(words, func(s string) bool { return len(s) > 10 }))
}

func TestContains2(t *testing.T) {
	// 测试键值对包含
	input := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
	assert.True(t, Contains2(input, 2, "b"))
	assert.False(t, Contains2(input, 3, "c"))
}

func TestIsSortedFunc(t *testing.T) {
	// 测试自定义排序判断
	numbers := fromSlice([]int{1, 3, 5, 7})
	assert.True(t, IsSortedFunc(numbers, cmp.Compare[int]))
	shuffled := fromSlice([]int{3, 1, 5, 7})
	assert.False(t, IsSortedFunc(shuffled, cmp.Compare[int]))
}

func TestCast(t *testing.T) {
	// 测试类型转换
	values := fromSlice([]any{1, "two", 3, "four"})
	casted := Cast[int](values)
	result := ToSlice2(casted)
	expected := []Pair[int, bool]{{1, true}, {0, false}, {3, true}, {0, false}}
	assert.Equal(t, expected, result)
}

func TestSizeValue(t *testing.T) {
	// 测试值计数
	numbers := fromSlice([]int{2, 5, 5, 7, 5})
	assert.Equal(t, 3, SizeValue(numbers, 5))
	assert.Equal(t, 0, SizeValue(numbers, 10))
}

func TestSizeValue2(t *testing.T) {
	// 测试键值对值计数
	pairs := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}, {3, "a"}})
	assert.Equal(t, 2, SizeValue2(pairs, "a"))
}

func TestFold2(t *testing.T) {
	// 测试键值对折叠
	input := FromSlice2([]Pair[int, int]{{1, 10}, {2, 20}, {3, 30}})
	sum := Fold2(input, 0, func(acc, k, v int) int { return acc + v })
	assert.Equal(t, 60, sum)
}

func TestAny2(t *testing.T) {
	// 测试键值对存在判断
	input := FromSlice2([]Pair[int, int]{{1, 5}, {2, 15}, {3, 25}})
	assert.True(t, Any2(input, func(k, v int) bool { return v > 20 }))
	assert.False(t, Any2(input, func(k, v int) bool { return v > 30 }))
}

func TestAll2(t *testing.T) {
	// 测试键值对全部满足判断
	input := FromSlice2([]Pair[int, int]{{1, 5}, {2, 15}, {3, 25}})
	assert.True(t, All2(input, func(k, v int) bool { return v > 0 }))
	assert.False(t, All2(input, func(k, v int) bool { return v < 20 }))
}

func TestEqualFunc2(t *testing.T) {
	// 测试键值对自定义比较
	seq1 := FromSlice2([]Pair[string, string]{{"a", "Apple"}, {"b", "Banana"}})
	seq2 := FromSlice2([]Pair[string, string]{{"a", "apple"}, {"b", "banana"}})
	assert.True(t, EqualFunc2(seq1, seq2, func(k1, v1, k2, v2 string) bool {
		return k1 == k2 && strings.EqualFold(v1, v2)
	}))
}

func TestMap2(t *testing.T) {
	// 测试键值对转换
	input := FromSlice2([]Pair[int, int]{{1, 10}, {2, 20}, {3, 30}})
	mapped := Map2(input, func(k, v int) (string, int) { return fmt.Sprintf("key%d", k), v * 2 })
	result := ToSlice2(mapped)
	assert.Equal(t, []Pair[string, int]{{"key1", 20}, {"key2", 40}, {"key3", 60}}, result)
}

func TestFilter2(t *testing.T) {
	// 测试键值对筛选
	input := FromSlice2([]Pair[int, int]{{1, 5}, {2, 15}, {3, 25}, {4, 5}})
	filtered := Filter2(input, func(k, v int) bool { return v > 10 })
	result := ToSlice2(filtered)
	assert.Equal(t, []Pair[int, int]{{2, 15}, {3, 25}}, result)
}

func TestContains(t *testing.T) {
	// 测试包含元素
	seq := fromSlice([]string{"a", "b", "c"})
	assert.True(t, Contains(seq, "b"))

	// 测试不包含元素
	assert.False(t, Contains(seq, "d"))
}

func TestIsSorted(t *testing.T) {
	// 测试已排序序列
	seq := fromSlice([]int{1, 2, 3, 4})
	assert.True(t, IsSorted(seq))

	// 测试未排序序列
	seq = fromSlice([]int{1, 3, 2, 4})
	assert.False(t, IsSorted(seq))
}

func TestZip(t *testing.T) {
    a := fromSlice([]int{1, 2, 3})
    b := fromSlice([]string{"a", "b"})
    z := Zip(a, b)
    got := ToSlice2(z)
    assert.Equal(t, []Pair[int, string]{{1, "a"}, {2, "b"}}, got)

    // 与空序列拉链
    empty := Empty[int]()
    got2 := ToSlice2(Zip(empty, fromSlice([]string{"x"})))
    assert.Empty(t, got2)
}

func TestZipWith(t *testing.T) {
    a := fromSlice([]int{1, 2, 3})
    b := fromSlice([]int{10, 20})
    z := ZipWith(a, b, func(x, y int) int { return x + y })
    got := ToSlice(z)
    assert.Equal(t, []int{11, 22}, got)
}

func TestZipWith2(t *testing.T) {
    a := FromSlice2([]Pair[int, string]{{1, "a"}, {2, "b"}})
    b := FromSlice2([]Pair[int, int]{{10, 3}})
    z := ZipWith2(a, b, func(k1 int, v1 string, k2 int, v2 int) (string, string) {
        return strconv.Itoa(k1+k2), v1+":"+strconv.Itoa(v2)
    })
    got := ToSlice2(z)
    assert.Equal(t, []Pair[string, string]{{"11", "a:3"}}, got)
}
