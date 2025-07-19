package xiter

import (
	"iter"
)

type integral interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type numbric interface {
	integral | ~float32 | ~float64
}

// Range1 生成从 0 到 end-1 的整数序列（不包含 end）。
//
// 示例: 生成 0 到 4 的序列
//
//	seq := xiter.Range1(5) // 生成序列: 0, 1, 2, 3, 4
//	result := xiter.ToSlice(seq)
//	// result 将为: [0 1 2 3 4]
func Range1[N integral](end N) iter.Seq[N] {
	return func(yield func(N) bool) {
		for i := N(0); i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Range2 生成从 start 到 end-1 的整数序列（不包含 end）。
//
// 示例: 生成 2 到 5 的序列
//
//	seq := xiter.Range2(2, 6) // 生成序列: 2, 3, 4, 5
//	result := xiter.ToSlice(seq)
//	// result 将为: [2 3 4 5]
func Range2[N integral](start, end N) iter.Seq[N] {
	return func(yield func(N) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Range3 生成从 start 到 end-1，步长为 step 的整数序列（不包含 end）。
//
// 示例: 生成 1 到 10 的奇数序列
//
//	seq := xiter.Range3(1, 10, 2) // 生成序列: 1, 3, 5, 7, 9
//	result := xiter.ToSlice(seq)
//	// result 将为: [1 3 5 7 9]
func Range3[N integral](start, end, step N) iter.Seq[N] {
	return func(yield func(N) bool) {
		for i := start; i < end; i += step {
			if !yield(i) {
				return
			}
		}
	}
}

// FromFunc 从函数生成序列，函数返回 (元素, 是否继续)，当返回 false 时序列结束。
//
// 示例: 从计数器函数生成序列
//
//	i := 0
//	seq := xiter.FromFunc(func() (int, bool) {
//		i++
//		return i, i <= 3
//	}) // 生成序列: 1, 2, 3
//	result := xiter.ToSlice(seq)
//	// result 将为: [1 2 3]
func FromFunc[E any](f func() (E, bool)) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			e, ok := f()
			if !ok {
				return
			}
			if !yield(e) {
				return
			}
		}
	}
}

// FromFunc2 从函数生成键值对序列，函数返回 (键, 值, 是否继续)，当返回 false 时序列结束。
//
// 示例: 从计数器函数生成键值对序列
//
//	i := 0
//	seq := xiter.FromFunc2(func() (int, string, bool) {
//		i++
//		return i, fmt.Sprintf("val%d", i), i <= 2
//	}) // 生成序列: (1,"val1"), (2,"val2")
//	result := xiter.ToSlice2(seq)
//	// result 将包含: [{Key:1 Value:"val1"}, {Key:2 Value:"val2"}]
func FromFunc2[K, V any](f func() (K, V, bool)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			k, v, ok := f()
			if !ok {
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// FromSlice 从切片生成键值对序列，键为切片索引，值为切片元素。
//
// 示例: 从字符串切片生成序列
//
//	slice := []string{"a", "b", "c"}
//	seq := xiter.FromSlice(slice) // 生成序列: (0,"a"), (1,"b"), (2,"c")
//	result := xiter.ToSlice2(seq)
//	// result 将包含: [{Key:0 Value:"a"}, {Key:1 Value:"b"}, {Key:2 Value:"c"}]
func FromSlice[E any, S ~[]E](s S) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i, e := range s {
			if !yield(i, e) {
				return
			}
		}
	}
}

// FromSlice2 从 Pair 类型切片生成键值对序列，每个 Pair 的 Key 和 Value 作为序列元素。
//
// 示例: 从 Pair 切片生成序列
//
//	slice := []xiter.Pair[int, string]{{1, "x"}, {2, "y"}}
//	seq := xiter.FromSlice2(slice) // 生成序列: (1,"x"), (2,"y")
//	result := xiter.ToSlice2(seq)
//	// result 将包含: [{Key:1 Value:"x"}, {Key:2 Value:"y"}]
func FromSlice2[K, V any, S ~[]Pair[K, V]](s S) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, p := range s {
			if !yield(p.Key, p.Value) {
				return
			}
		}
	}
}

// FromMap 将映射转换为键值对序列，遍历映射中的所有键值对。
//
// 示例: 从映射生成序列
//
//	m := map[string]int{"a": 1, "b": 2}
//	seq := xiter.FromMap(m) // 生成序列: ("a",1), ("b",2)
//	result := xiter.ToSlice2(seq)
//	// result 将包含: [{Key:"a" Value:1}, {Key:"b" Value:2}]
func FromMap[K comparable, V any, M ~map[K]V](m M) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Once 生成只包含单个元素的序列。
//
// 示例: 生成只包含一个元素的序列
//
//	seq := xiter.Once("hello") // 生成序列: "hello"
//	result := xiter.ToSlice(seq)
//	// result 将为: [hello]
func Once[E any](e E) iter.Seq[E] {
	return func(yield func(E) bool) {
		yield(e)
	}
}

// Once2 生成只包含单个键值对的序列。
//
// 示例: 生成只包含一个键值对的序列
//
//	seq := xiter.Once2(1, "x") // 生成序列: (1,"x")
//	result := xiter.ToSlice2(seq)
//	// result 将包含: [{Key:1 Value:"x"}]
func Once2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(k, v)
	}
}

// Empty 生成不包含任何元素的空序列。
//
// 示例: 生成空序列
//
//	seq := xiter.Empty[int]() // 生成空序列
//	result := xiter.ToSlice(seq)
//	// result 将为: []
func Empty[E any]() iter.Seq[E] {
	return func(yield func(E) bool) {
	}
}

// Empty2 生成不包含任何键值对的空序列。
//
// 示例: 生成空键值对序列
//
//	seq := xiter.Empty2[int, string]() // 生成空序列
//	result := xiter.ToSlice2(seq)
//	// result 将为: []
func Empty2[K, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
	}
}

// Repeat 生成无限重复单个元素的序列。
//
// 示例: 生成重复 "a" 的序列（通常与 Take 配合使用限制长度）
//
//	seq := xiter.Repeat("a")
//	limited := xiter.Take(seq, 3) // 取前3个元素
//	result := xiter.ToSlice(limited)
//	// result 将为: [a a a]
func Repeat[E any](e E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			if !yield(e) {
				return
			}
		}
	}
}

// Repeat2 生成无限重复单个键值对的序列。
//
// 示例: 生成重复 (1,"x") 的序列（通常与 Take 配合使用限制长度）
//
//	seq := xiter.Repeat2(1, "x")
//	limited := xiter.Take2(seq, 2) // 取前2个键值对
//	result := xiter.ToSlice2(limited)
//	// result 将包含: [{Key:1 Value:"x"}, {Key:1 Value:"x"}]
func Repeat2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			if !yield(k, v) {
				return
			}
		}
	}
}
