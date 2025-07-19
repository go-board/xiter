package xiter

import (
	"iter"
)

type Pair[K, V any] struct {
	Key   K
	Value V
}

// GroupBy 根据指定的键函数对序列元素进行分组，返回 map[K][]E 类型的分组结果。
//
// 示例: 将字符串按首字母分组
//
//	words := xiter.FromSlice([]string{"apple", "banana", "apricot", "cherry"})
//	groups := xiter.GroupBy(words, func(s string) rune { return []rune(s)[0] })
//	// groups 将为: map[a:[apple apricot] b:[banana] c:[cherry]]
func GroupBy[E any, K comparable](s iter.Seq[E], f func(E) K) map[K][]E {
	m := make(map[K][]E)
	for e := range s {
		k := f(e)
		m[k] = append(m[k], e)
	}
	return m
}

// ToSet 将序列转换为 map[E]struct{} 类型的集合，返回包含序列所有唯一元素的新集合。
//
// 示例: 将整数序列转换为集合
//
//	seq := xiter.FromSlice([]int{1, 2, 2, 3})
//	result := xiter.ToSet(seq) // result 类型为 map[int]struct{}
//	// result 将包含元素: 1, 2, 3
func ToSet[E comparable](s iter.Seq[E]) map[E]struct{} {
	set := make(map[E]struct{})
	for e := range s {
		set[e] = struct{}{}
	}
	return set
}

// ToSlice 将序列转换为 []E 类型的切片，返回包含序列所有元素的新切片。
//
// 示例 1: 将整数序列转换为 []int 切片
//
//	seq := xiter.FromSlice([]int{1, 2, 3, 4})
//	result := xiter.ToSlice(seq) // result 类型为 []int
//	// result 将为: [1 2 3 4]
//
// 示例 2: 将字符串序列转换为 []string 切片
//
//	words := xiter.FromSlice([]string{"a", "b", "c"})
//	strs := xiter.ToSlice(words) // strs 类型为 []string
//	// strs 将为: [a b c]
func ToSlice[E any](s iter.Seq[E]) []E {
	var slice []E
	for e := range s {
		slice = append(slice, e)
	}
	return slice
}

// ToMap 将键值对序列转换为 map[K]V 类型的映射，返回包含所有键值对的新映射。
//
// 示例: 将键值对序列转换为映射
//
//	seq := xiter.FromSlice2([]iter.Pair[string, int]{{"a", 1}, {"b", 2}})
//	result := xiter.ToMap(seq) // result 类型为 map[string]int
//	// result 将为: map[a:1 b:2]
func ToMap[K comparable, V any](s iter.Seq2[K, V]) map[K]V {
	m := make(map[K]V)
	for k, v := range s {
		m[k] = v
	}
	return m
}

// ToSlice2 将键值对序列转换为 []Pair[K, V] 类型的切片，返回包含所有键值对的新切片。
//
// 示例: 将键值对序列转换为切片
//
//	seq := xiter.FromSlice2([]iter.Pair[int, string]{{1, "a"}, {2, "b"}})
//	result := xiter.ToSlice2(seq) // result 类型为 []Pair[int, string]
//	// result 将为: [{Key:1 Value:"a"} {Key:2 Value:"b"}]
func ToSlice2[K, V any](s iter.Seq2[K, V]) []Pair[K, V] {
	var slice []Pair[K, V]
	for k, v := range s {
		slice = append(slice, Pair[K, V]{Key: k, Value: v})
	}
	return slice
}
