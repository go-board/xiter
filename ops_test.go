package xiter

import (
	"reflect"
	"strings"
	"testing"
)

func TestTake(t *testing.T) {
	// 创建一个测试序列
	s := Range1(10)

	// 使用Take获取前5个元素
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

	// 测试获取超过序列长度的元素
	s = Range1(5)
	taken = Take(s, 10)
	result = ToSlice(taken)
	expected = []int{0, 1, 2, 3, 4}

	if len(result) != len(expected) {
		t.Errorf("Take returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试获取0个元素
	taken = Take(s, 0)
	result = ToSlice(taken)
	if len(result) != 0 {
		t.Errorf("Take returned %d elements, expected 0", len(result))
	}

	// 测试空序列
	emptySeq := Empty[int]()
	taken = Take(emptySeq, 5)
	result = ToSlice(taken)
	if len(result) != 0 {
		t.Errorf("Take on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestTake2(t *testing.T) {
	// 创建一个测试键值对序列
	values := []int{10, 20, 30, 40, 50}
	s := FromSlice(values)

	// 使用Take2获取前3个元素
	taken := Take2(s, 3)
	result := ToMap(taken)
	expected := map[int]int{0: 10, 1: 20, 2: 30}

	if len(result) != len(expected) {
		t.Errorf("Take2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Take2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	// 测试获取超过序列长度的元素
	taken = Take2(s, 10)
	result = ToMap(taken)
	expected = map[int]int{0: 10, 1: 20, 2: 30, 3: 40, 4: 50}

	if len(result) != len(expected) {
		t.Errorf("Take2 returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试获取0个元素
	taken = Take2(s, 0)
	result = ToMap(taken)
	if len(result) != 0 {
		t.Errorf("Take2 returned %d elements, expected 0", len(result))
	}

	// 测试空序列
	emptySeq := Empty2[int, int]()
	taken = Take2(emptySeq, 5)
	result = ToMap(taken)
	if len(result) != 0 {
		t.Errorf("Take2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestSkip(t *testing.T) {
	// 创建一个测试序列
	s := Range1(10)

	// 使用Skip跳过前5个元素
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

	// 测试跳过超过序列长度的元素
	s = Range1(5)
	skiped = Skip(s, 10)
	result = ToSlice(skiped)
	if len(result) != 0 {
		t.Errorf("Skip returned %d elements, expected 0", len(result))
	}

	// 测试跳过0个元素
	skiped = Skip(s, 0)
	result = ToSlice(skiped)
	expected = []int{0, 1, 2, 3, 4}

	if len(result) != len(expected) {
		t.Errorf("Skip returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试空序列
	emptySeq := Empty[int]()
	skiped = Skip(emptySeq, 5)
	result = ToSlice(skiped)
	if len(result) != 0 {
		t.Errorf("Skip on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestSkip2(t *testing.T) {
	// 创建一个测试键值对序列
	values := []int{10, 20, 30, 40, 50}
	s := FromSlice(values)

	// 使用Skip2跳过前3个元素
	skiped := Skip2(s, 3)
	result := ToMap(skiped)
	expected := map[int]int{3: 40, 4: 50}

	if len(result) != len(expected) {
		t.Errorf("Skip2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Skip2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	// 测试跳过超过序列长度的元素
	skiped = Skip2(s, 10)
	result = ToMap(skiped)
	if len(result) != 0 {
		t.Errorf("Skip2 returned %d elements, expected 0", len(result))
	}

	// 测试跳过0个元素
	skiped = Skip2(s, 0)
	result = ToMap(skiped)
	expected = map[int]int{0: 10, 1: 20, 2: 30, 3: 40, 4: 50}

	if len(result) != len(expected) {
		t.Errorf("Skip2 returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试空序列
	emptySeq := Empty2[int, int]()
	skiped = Skip2(emptySeq, 5)
	result = ToMap(skiped)
	if len(result) != 0 {
		t.Errorf("Skip2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestEnumerate(t *testing.T) {
	// 创建一个测试序列
	values := []string{"a", "b", "c", "d", "e"}
	s := func(yield func(string) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}

	// 使用Enumerate为每个元素添加索引
	enumerated := Enumerate(s)
	result := ToMap(enumerated)

	// 验证结果
	expected := map[int]string{0: "a", 1: "b", 2: "c", 3: "d", 4: "e"}
	if len(result) != len(expected) {
		t.Errorf("Enumerate returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Enumerate[%d] = %s, expected %s", k, result[k], v)
		}
	}

	// 测试空序列
	emptySeq := Empty[string]()
	enumerated = Enumerate(emptySeq)
	result = ToMap(enumerated)
	if len(result) != 0 {
		t.Errorf("Enumerate on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestMap(t *testing.T) {
	// 创建一个测试序列
	s := Range1(5)

	// 使用Map将每个元素乘以2
	mapped := Map(s, func(x int) int {
		return x * 2
	})
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

	// 测试类型转换
	mappedStr := Map(s, func(x int) string {
		return string(rune('a' + x))
	})
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

	// 测试空序列
	emptySeq := Empty[int]()
	mapped = Map(emptySeq, func(x int) int {
		return x * 2
	})
	result = ToSlice(mapped)
	if len(result) != 0 {
		t.Errorf("Map on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestMap2(t *testing.T) {
	// 创建测试序列
	values := map[string]int{"a": 1, "b": 2, "c": 3}
	s := FromMap(values)

	// 使用Map2转换
	mapped := Map2(s, func(k string, v int) (int, string) {
		return v, k
	})

	// 验证结果
	result := ToMap(mapped)
	if len(result) != 3 {
		t.Errorf("TestMap2: expected 3 elements, got %d", len(result))
	}
	expected := map[int]string{1: "a", 2: "b", 3: "c"}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("TestMap2: expected value '%s' for key %d, got '%s'", v, k, result[k])
		}
	}

	// 测试空序列
	emptySeq := Empty2[string, int]()
	emptyResult := Map2(emptySeq, func(k string, v int) (int, string) {
		return v, k
	})
	emptyMap := ToMap(emptyResult)
	if len(emptyMap) != 0 {
		t.Errorf("TestMap2: expected empty map for empty sequence, got %d elements", len(emptyMap))
	}
}

func TestFilter(t *testing.T) {
	// 创建一个测试序列
	s := Range1(10)

	// 使用Filter过滤出偶数
	filtered := Filter(s, func(x int) bool {
		return x%2 == 0
	})
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

	// 测试过滤出所有元素
	filtered = Filter(s, func(x int) bool {
		return true
	})
	result = ToSlice(filtered)
	expected = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	if len(result) != len(expected) {
		t.Errorf("Filter returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试过滤出0个元素
	filtered = Filter(s, func(x int) bool {
		return false
	})
	result = ToSlice(filtered)
	if len(result) != 0 {
		t.Errorf("Filter returned %d elements, expected 0", len(result))
	}

	// 测试空序列
	emptySeq := Empty[int]()
	filtered = Filter(emptySeq, func(x int) bool {
		return x%2 == 0
	})
	result = ToSlice(filtered)
	if len(result) != 0 {
		t.Errorf("Filter on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestFilter2(t *testing.T) {
	// 创建一个测试键值对序列
	values := []int{10, 20, 30, 40, 50}
	s := FromSlice(values)

	// 使用Filter2过滤出键为偶数的元素
	filtered := Filter2(s, func(k int, v int) bool {
		return k%2 == 0
	})
	result := ToMap(filtered)
	expected := map[int]int{0: 10, 2: 30, 4: 50}

	if len(result) != len(expected) {
		t.Errorf("Filter2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Filter2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	// 测试过滤出所有元素
	filtered = Filter2(s, func(k int, v int) bool {
		return true
	})
	result = ToMap(filtered)
	expected = map[int]int{0: 10, 1: 20, 2: 30, 3: 40, 4: 50}

	if len(result) != len(expected) {
		t.Errorf("Filter2 returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试过滤出0个元素
	filtered = Filter2(s, func(k int, v int) bool {
		return false
	})
	result = ToMap(filtered)
	if len(result) != 0 {
		t.Errorf("Filter2 returned %d elements, expected 0", len(result))
	}

	// 测试空序列
	emptySeq := Empty2[int, int]()
	filtered = Filter2(emptySeq, func(k int, v int) bool {
		return true
	})
	result = ToMap(filtered)
	if len(result) != 0 {
		t.Errorf("Filter2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestTakeWhile(t *testing.T) {
	// 创建一个测试序列
	s := Range1(10)

	// 使用TakeWhile获取元素直到遇到大于5的元素
	takenWhile := TakeWhile(s, func(x int) bool {
		return x <= 5
	})
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

	// 测试条件始终为假
	takenWhile = TakeWhile(s, func(x int) bool {
		return false
	})
	result = ToSlice(takenWhile)
	if len(result) != 0 {
		t.Errorf("TakeWhile returned %d elements, expected 0", len(result))
	}

	// 测试条件始终为真
	takenWhile = TakeWhile(s, func(x int) bool {
		return true
	})
	result = ToSlice(takenWhile)
	expected = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	if len(result) != len(expected) {
		t.Errorf("TakeWhile returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试空序列
	emptySeq := Empty[int]()
	takenWhile = TakeWhile(emptySeq, func(x int) bool {
		return true
	})
	result = ToSlice(takenWhile)
	if len(result) != 0 {
		t.Errorf("TakeWhile on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestTakeWhile2(t *testing.T) {
	// 创建一个测试键值对序列
	values := []int{10, 20, 30, 40, 50}
	s := FromSlice(values)

	// 使用TakeWhile2获取元素直到遇到键大于2的元素
	takenWhile := TakeWhile2(s, func(k int, v int) bool {
		return k <= 2
	})
	result := ToMap(takenWhile)
	expected := map[int]int{0: 10, 1: 20, 2: 30}

	if len(result) != len(expected) {
		t.Errorf("TakeWhile2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("TakeWhile2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	// 测试条件始终为假
	takenWhile = TakeWhile2(s, func(k int, v int) bool {
		return false
	})
	result = ToMap(takenWhile)
	if len(result) != 0 {
		t.Errorf("TakeWhile2 returned %d elements, expected 0", len(result))
	}

	// 测试条件始终为真
	takenWhile = TakeWhile2(s, func(k int, v int) bool {
		return true
	})
	result = ToMap(takenWhile)
	expected = map[int]int{0: 10, 1: 20, 2: 30, 3: 40, 4: 50}

	if len(result) != len(expected) {
		t.Errorf("TakeWhile2 returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试空序列
	emptySeq := Empty2[int, int]()
	takenWhile = TakeWhile2(emptySeq, func(k int, v int) bool {
		return true
	})
	result = ToMap(takenWhile)
	if len(result) != 0 {
		t.Errorf("TakeWhile2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestSkipWhile(t *testing.T) {
	// 创建一个测试序列
	s := Range1(10)

	// 使用SkipWhile跳过元素直到遇到大于5的元素
	skipedWhile := SkipWhile(s, func(x int) bool {
		return x <= 5
	})
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

	// 测试条件始终为假
	skipedWhile = SkipWhile(s, func(x int) bool {
		return false
	})
	result = ToSlice(skipedWhile)
	expected = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	if len(result) != len(expected) {
		t.Errorf("SkipWhile returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试条件始终为真
	skipedWhile = SkipWhile(s, func(x int) bool {
		return true
	})
	result = ToSlice(skipedWhile)
	if len(result) != 0 {
		t.Errorf("SkipWhile returned %d elements, expected 0", len(result))
	}

	// 测试空序列
	emptySeq := Empty[int]()
	skipedWhile = SkipWhile(emptySeq, func(x int) bool {
		return true
	})
	result = ToSlice(skipedWhile)
	if len(result) != 0 {
		t.Errorf("SkipWhile on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestSkipWhile2(t *testing.T) {
	// 创建一个测试键值对序列
	values := []int{10, 20, 30, 40, 50}
	s := FromSlice(values)

	// 使用SkipWhile2跳过元素直到遇到键大于2的元素
	skipedWhile := SkipWhile2(s, func(k int, v int) bool {
		return k <= 2
	})
	result := ToMap(skipedWhile)
	expected := map[int]int{3: 40, 4: 50}

	if len(result) != len(expected) {
		t.Errorf("SkipWhile2 returned %d elements, expected %d", len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("SkipWhile2[%d] = %d, expected %d", k, result[k], v)
		}
	}

	// 测试条件始终为假
	skipedWhile = SkipWhile2(s, func(k int, v int) bool {
		return false
	})
	result = ToMap(skipedWhile)
	expected = map[int]int{0: 10, 1: 20, 2: 30, 3: 40, 4: 50}

	if len(result) != len(expected) {
		t.Errorf("SkipWhile2 returned %d elements, expected %d", len(result), len(expected))
	}

	// 测试条件始终为真
	skipedWhile = SkipWhile2(s, func(k int, v int) bool {
		return true
	})
	result = ToMap(skipedWhile)
	if len(result) != 0 {
		t.Errorf("SkipWhile2 returned %d elements, expected 0", len(result))
	}

	// 测试空序列
	emptySeq := Empty2[int, int]()
	skipedWhile = SkipWhile2(emptySeq, func(k int, v int) bool {
		return true
	})
	result = ToMap(skipedWhile)
	if len(result) != 0 {
		t.Errorf("SkipWhile2 on empty sequence returned %d elements, expected 0", len(result))
	}
}

func TestFold(t *testing.T) {
	// 创建一个测试序列
	s := Range1(5)

	// 使用Fold计算元素总和
	sum := Fold(s, 0, func(acc, e int) int {
		return acc + e
	})
	expectedSum := 10 // 0 + 1 + 2 + 3 + 4
	if sum != expectedSum {
		t.Errorf("Fold sum: expected %d, got %d", expectedSum, sum)
	}

	// 使用Fold连接字符串
	sStr := func(yield func(string) bool) {
		yield("a")
		yield("b")
		yield("c")
	}
	concatenated := Fold(sStr, "", func(acc, e string) string {
		return acc + e
	})
	expectedConcatenated := "abc"
	if concatenated != expectedConcatenated {
		t.Errorf("Fold concatenation: expected %s, got %s", expectedConcatenated, concatenated)
	}

	// 测试空序列
	emptySeq := Empty[int]()
	emptyResult := Fold(emptySeq, 42, func(acc, e int) int {
		return acc + e
	})
	if emptyResult != 42 {
		t.Errorf("Fold on empty sequence: expected 42, got %d", emptyResult)
	}
}

func TestFold2(t *testing.T) {
	// 创建一个测试键值对序列
	values := map[string]int{"a": 1, "b": 2, "c": 3}
	s := FromMap(values)

	// 使用Fold2计算值的总和
	sum := Fold2(s, 0, func(acc int, k string, v int) int {
		return acc + v
	})
	expectedSum := 6 // 1 + 2 + 3
	if sum != expectedSum {
		t.Errorf("Fold2 sum: expected %d, got %d", expectedSum, sum)
	}

	// 使用Fold2构建一个新的map
	newMap := Fold2(s, make(map[int]string), func(acc map[int]string, k string, v int) map[int]string {
		acc[v] = k
		return acc
	})
	expectedMap := map[int]string{1: "a", 2: "b", 3: "c"}
	if len(newMap) != len(expectedMap) {
		t.Errorf("Fold2 map construction: expected %d elements, got %d", len(expectedMap), len(newMap))
	}
	for k, v := range expectedMap {
		if newMap[k] != v {
			t.Errorf("Fold2 map construction: expected %s for key %d, got %s", v, k, newMap[k])
		}
	}

	// 测试空序列
	emptySeq := Empty2[string, int]()
	emptyResult := Fold2(emptySeq, "initial", func(acc string, k string, v int) string {
		return acc + k
	})
	if emptyResult != "initial" {
		t.Errorf("Fold2 on empty sequence: expected 'initial', got %s", emptyResult)
	}
}

func TestSize(t *testing.T) {
	// 测试正常序列
	s := Range1(5)
	size := Size(s)
	expected := 5
	if size != expected {
		t.Errorf("Size: expected %d, got %d", expected, size)
	}

	// 测试空序列
	emptySeq := Empty[int]()
	emptySize := Size(emptySeq)
	if emptySize != 0 {
		t.Errorf("Size on empty sequence: expected 0, got %d", emptySize)
	}

	// 测试有限的无限序列
	infinite := Repeat("test")
	limited := Take(infinite, 10)
	limitedSize := Size(limited)
	if limitedSize != 10 {
		t.Errorf("Size on limited infinite sequence: expected 10, got %d", limitedSize)
	}
}

func TestSize2(t *testing.T) {
	// 测试正常键值对序列
	values := map[string]int{"a": 1, "b": 2, "c": 3}
	s := FromMap(values)
	size := Size2(s)
	expected := 3
	if size != expected {
		t.Errorf("Size2: expected %d, got %d", expected, size)
	}

	// 测试空序列
	emptySeq := Empty2[string, int]()
	emptySize := Size2(emptySeq)
	if emptySize != 0 {
		t.Errorf("Size2 on empty sequence: expected 0, got %d", emptySize)
	}

	// 测试有限的无限序列
	infinite := Repeat2("key", "value")
	limited := Take2(infinite, 5)
	limitedSize := Size2(limited)
	if limitedSize != 5 {
		t.Errorf("Size2 on limited infinite sequence: expected 5, got %d", limitedSize)
	}
}

func TestSizeFunc(t *testing.T) {
	// 测试筛选偶数
	s := Range1(10)
	evenCount := SizeFunc(s, func(e int) bool {
		return e%2 == 0
	})
	expected := 5 // 0, 2, 4, 6, 8
	if evenCount != expected {
		t.Errorf("SizeFunc (even numbers): expected %d, got %d", expected, evenCount)
	}

	// 测试筛选大于10的数（应该没有）
	smallerCount := SizeFunc(s, func(e int) bool {
		return e > 10
	})
	if smallerCount != 0 {
		t.Errorf("SizeFunc (numbers > 10): expected 0, got %d", smallerCount)
	}

	// 测试空序列
	emptySeq := Empty[int]()
	emptyCount := SizeFunc(emptySeq, func(e int) bool {
		return true
	})
	if emptyCount != 0 {
		t.Errorf("SizeFunc on empty sequence: expected 0, got %d", emptyCount)
	}
}

// 修复拼写错误
func TestSizeFunc2(t *testing.T) {
	// 测试筛选值大于1的键值对
	values := map[string]int{"a": 1, "b": 2, "c": 3, "d": 1}
	s := FromMap(values)
	largerThanOneCount := SizeFunc2(s, func(k string, v int) bool {
		return v > 1
	})
	expected := 2 // "b":2, "c":3
	if largerThanOneCount != expected {
		t.Errorf("SizeFunc2 (values > 1): expected %d, got %d", expected, largerThanOneCount)
	}

	// 测试筛选键包含特定字符
	keyWithACount := SizeFunc2(s, func(k string, v int) bool {
		return strings.Contains(k, "a")
	})
	expected = 1 // only "a"
	if keyWithACount != expected {
		t.Errorf("SizeFunc2 (keys with 'a'): expected %d, got %d", expected, keyWithACount)
	}

	// 测试空序列
	emptySeq := Empty2[string, int]()
	emptyCount := SizeFunc2(emptySeq, func(k string, v int) bool {
		return true
	})
	if emptyCount != 0 {
		t.Errorf("SizeFunc2 on empty sequence: expected 0, got %d", emptyCount)
	}
}

// TestSizeValue2 tests SizeValue2 function
func TestSizeValue2(t *testing.T) {
	// 测试匹配特定键值对
	s := func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
		yield("c", 3)
		yield("a", 1) // 添加一个重复的键值对
	}
	expected := 2 // 有两个 "a": 1
	count := SizeValue2(s, "a", 1)
	if count != expected {
		t.Errorf("SizeValue2 (key='a', value=1): expected %d, got %d", expected, count)
	}

	// 测试匹配不存在的键值对
	count = SizeValue2(s, "d", 4)
	if count != 0 {
		t.Errorf("SizeValue2 (key='d', value=4): expected 0, got %d", count)
	}

	// 测试只匹配键或值
	count = SizeValue2(s, "a", 2) // 键匹配但值不匹配
	if count != 0 {
		t.Errorf("SizeValue2 (key='a', value=2): expected 0, got %d", count)
	}

	count = SizeValue2(s, "b", 1) // 值匹配但键不匹配
	if count != 0 {
		t.Errorf("SizeValue2 (key='b', value=1): expected 0, got %d", count)
	}

	// 测试空序列
	emptySeq := Empty2[string, int]()
	count = SizeValue2(emptySeq, "a", 1)
	if count != 0 {
		t.Errorf("SizeValue2 on empty sequence: expected 0, got %d", count)
	}
}

// 添加Compare系列函数的测试用例
func TestCompare(t *testing.T) {
	// 测试相等序列
	s1 := Range1(3)
	s2 := Range1(3)
	result := Compare(s1, s2)
	if result != 0 {
		t.Errorf("Compare (equal sequences): expected 0, got %d", result)
	}

	// 测试s1 < s2
	s3 := Range1(2)
	s4 := Range1(3)
	result = Compare(s3, s4)
	if result >= 0 {
		t.Errorf("Compare (s1 < s2): expected < 0, got %d", result)
	}

	// 测试s1 > s2
	s5 := Range1(3)
	s6 := Range1(2)
	result = Compare(s5, s6)
	if result <= 0 {
		t.Errorf("Compare (s1 > s2): expected > 0, got %d", result)
	}

	// 测试元素值不同的序列
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
	// 测试自定义比较函数
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

	// 测试自定义比较（例如绝对值比较）
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

func TestCompare2(t *testing.T) {
	// 测试相等键值对序列
	s1 := func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
	}
	s2 := func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
	}
	result := Compare2(s1, s2)
	if result != 0 {
		t.Errorf("Compare2 (equal sequences): expected 0, got %d", result)
	}

	// 测试键值对值不同的序列
	s3 := func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
	}
	s4 := func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 3)
	}
	result = Compare2(s3, s4)
	if result >= 0 {
		t.Errorf("Compare2 (values differ): expected < 0, got %d", result)
	}
}

func TestCompareFunc2(t *testing.T) {
	// 测试自定义键值对比较函数
	s1 := func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
	}
	s2 := func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
	}
	result := CompareFunc2(s1, s2, func(k1 string, v1 int, k2 string, v2 int) int {
		if k1 < k2 {
			return -1
		} else if k1 > k2 {
			return 1
		}
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	})
	if result != 0 {
		t.Errorf("CompareFunc2 (equal sequences): expected 0, got %d", result)
	}

	// 测试自定义比较（只比较值）
	s3 := func(yield func(string, int) bool) {
		yield("x", 1)
		yield("y", 2)
	}
	s4 := func(yield func(string, int) bool) {
		yield("a", 1)
		yield("b", 2)
	}
	result = CompareFunc2(s3, s4, func(k1 string, v1 int, k2 string, v2 int) int {
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	})
	if result != 0 {
		t.Errorf("CompareFunc2 (values equal): expected 0, got %d", result)
	}
}

// TestFilterMap verifies that FilterMap skips elements when the function returns false
func TestFilterMap(t *testing.T) {
	// Create a sequence: 1, 2, 3, 4, 5
	source := func(yield func(int) bool) {
		for i := 1; i <= 5; i++ {
			if !yield(i) {
				return
			}
		}
	}

	// Define FilterMap function: only keep even numbers and double them
	// For even elements: return (element*2, true)
	// For odd elements: return (0, false) - should skip
	filterFn := func(x int) (int, bool) {
		if x%2 == 0 {
			return x * 2, true
		}
		return 0, false
	}

	// Apply FilterMap
	result := FilterMap(source, filterFn)

	// Expected results: 4 (2*2), 8 (4*2)
	expected := []int{4, 8}
	var actual []int
	for x := range result {
		actual = append(actual, x)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FilterMap failed. Expected: %v, Got: %v", expected, actual)
	}
}

// TestFilterMap2 verifies that FilterMap2 skips elements when the function returns false
func TestFilterMap2(t *testing.T) {
	// Create a sequence: ("a", 1), ("b", 2), ("c", 3), ("d", 4)
	source := func(yield func(string, int) bool) {
		for i, c := range []string{"a", "b", "c", "d"} {
			if !yield(c, i+1) {
				return
			}
		}
	}

	// Define FilterMap2 function: only keep pairs where value is even
	// For even values: return (key+key, value*2, true)
	// For odd values: return ("", 0, false) - should skip
	filterFn := func(k string, v int) (string, int, bool) {
		if v%2 == 0 {
			return k + k, v * 2, true
		}
		return "", 0, false
	}

	// Apply FilterMap2
	result := FilterMap2(source, filterFn)

	// Expected results: ("bb", 4), ("dd", 8)
	expected := map[string]int{"bb": 4, "dd": 8}
	var actualKeys []string
	var actualValues []int
	for k, v := range result {
		actualKeys = append(actualKeys, k)
		actualValues = append(actualValues, v)
		if expected[k] != v {
			t.Errorf("FilterMap2 failed for key %s. Expected: %d, Got: %d", k, expected[k], v)
		}
		delete(expected, k)
	}

	if len(expected) > 0 {
		t.Errorf("FilterMap2 failed. Missing keys: %v", expected)
	}

	if len(actualKeys) != 2 {
		t.Errorf("FilterMap2 failed. Expected 2 elements, Got: %d", len(actualKeys))
	}
}
