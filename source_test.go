package xiter

import (
	"testing"
)

func TestRange1(t *testing.T) {
	// 测试正常范围
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

	// 测试边界情况：0
	s = Range1(0)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range1(0) returned %d elements, expected 0", len(result))
	}

	// 测试边界情况：1
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
	// 测试正常范围
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

	// 测试边界情况：start == end
	s = Range2(5, 5)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range2(5,5) returned %d elements, expected 0", len(result))
	}

	// 测试边界情况：start > end
	s = Range2(7, 2)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range2(7,2) returned %d elements, expected 0", len(result))
	}
}

func TestRange3(t *testing.T) {
	// 测试正常范围和步长
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

	// 测试边界情况：步长为0（应该不会进入循环）
	s = Range3(1, 10, 0)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range3(1,10,0) returned %d elements, expected 0", len(result))
	}

	// 测试边界情况：start >= end
	s = Range3(10, 1, 2)
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Range3(10,1,2) returned %d elements, expected 0", len(result))
	}
}

func TestFromFunc(t *testing.T) {
	// 测试正常功能
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

	// 测试空序列
	s = FromFunc(func() (int, bool) {
		return 0, false
	})
	result = ToSlice(s)
	if len(result) != 0 {
		t.Errorf("FromFunc returned %d elements, expected 0", len(result))
	}
}

func TestFromFunc2(t *testing.T) {
	// 测试正常功能
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

	// 测试空序列
	s = FromFunc2(func() (string, int, bool) {
		return "", 0, false
	})
	result = ToMap(s)
	if len(result) != 0 {
		t.Errorf("FromFunc2 returned %d elements, expected 0", len(result))
	}
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

func TestEmpty(t *testing.T) {
	s := Empty[int]()
	result := ToSlice(s)
	if len(result) != 0 {
		t.Errorf("Empty returned %d elements, expected 0", len(result))
	}
}

func TestEmpty2(t *testing.T) {
	s := Empty2[string, int]()
	result := ToMap(s)
	if len(result) != 0 {
		t.Errorf("Empty2 returned %d elements, expected 0", len(result))
	}
}

func TestRepeat(t *testing.T) {
	s := Repeat("test")
	// 使用Take来限制获取的元素数量，因为Repeat是无限的
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

func TestRepeat2(t *testing.T) {
	s := Repeat2("key", "value")
	// 使用Take2来限制获取的元素数量，因为Repeat2是无限的
	taken := Take2(s, 3)
	result := ToMap(taken)
	// 注意：由于是重复的键值对，最终map中只会有一个元素
	if len(result) != 1 {
		t.Errorf("Repeat2 returned %d elements, expected 1", len(result))
	}
	if result["key"] != "value" {
		t.Errorf("Repeat2 returned %s, expected %s", result["key"], "value")
	}
}
