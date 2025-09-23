package xiter

import (
	"cmp"
	"iter"
)

// Take 返回序列的前 n 个元素
func Take[E any](s iter.Seq[E], n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		count := 0
		for e := range s {
			if count >= n {
				return
			}
			if !yield(e) {
				return
			}
			count++
		}
	}
}

// Take2 返回键值对序列的前 n 个元素
func Take2[K, V any](s iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		count := 0
		for k, v := range s {
			if count >= n {
				return
			}
			if !yield(k, v) {
				return
			}
			count++
		}
	}
}

// Enumerate 返回一个包含 (索引, 元素) 元组的序列，为输入序列中的每个元素添加索引。
//
// 示例 1: 为数字序列添加索引
//
//	seq := xiter.Range(0, 5)       // 生成 0, 1, 2, 3, 4
//	enumerated := xiter.Enumerate(seq)
//	xiter.ForEach2(enumerated, func(i int, e int) {
//		fmt.Printf("(%d, %d) ", i, e)
//	})
//
// 输出:
//
//	(0, 0) (1, 1) (2, 2) (3, 3) (4, 4)
//
// 示例 2: 为字符串序列添加索引
//
//	words := xiter.FromSlice([]string{"apple", "banana", "cherry"})
//	enumerated := xiter.Enumerate(words)
//	result, _ := xiter.ToSlice2(enumerated)
//	// result 将包含: [(0, "apple"), (1, "banana"), (2, "cherry")]
func Enumerate[E any](s iter.Seq[E]) iter.Seq2[int, E] {
	i := -1
	return func(yield func(int, E) bool) {
		for e := range s {
			i++
			if !yield(i, e) {
				return
			}
		}
	}
}

// Map 对输入序列中的每个元素应用指定函数，并返回包含结果的新序列。
//
// 示例 1: 将整数序列中的每个元素乘以 2
//
//	input := xiter.Range(0, 5)      // 生成序列: 0, 1, 2, 3, 4
//	doubled := xiter.Map(input, func(e int) int { return e * 2 })
//	result, _ := xiter.ToSlice(doubled)
//	// result 将包含: [0, 2, 4, 6, 8]
//
// 示例 2: 将整数转换为字符串
//
//	numbers := xiter.Range(1, 4)     // 生成序列: 1, 2, 3
//	strings := xiter.Map(numbers, func(n int) string { return fmt.Sprintf("%d", n) })
//	result, _ := xiter.ToSlice(strings)
//	// result 将包含: ["1", "2", "3"]
func Map[E1, E2 any](s iter.Seq[E1], f func(E1) E2) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e := range s {
			if !yield(f(e)) {
				return
			}
		}
	}
}

// Map2 对输入键值对序列中的每个元素应用指定函数，并返回包含结果的新键值对序列。
//
// 示例: 将键值对中的值加倍
//
//	input := xiter.Enumerate(xiter.Range(1, 4)) // 生成序列: (0,1), (1,2), (2,3)
//	doubled := xiter.Map2(input, func(k int, v int) (int, int) { return k, v * 2 })
//	result, _ := xiter.ToSlice2(doubled)
//	// result 将包含: [(0, 2), (1, 4), (2, 6)]
func Map2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Filter 返回一个仅包含满足谓词函数条件的元素的新序列。
//
// 示例 1: 筛选偶数
//
//	input := xiter.Range(0, 10)     // 生成序列: 0-9
//	evens := xiter.Filter(input, func(e int) bool { return e%2 == 0 })
//	result, _ := xiter.ToSlice(evens)
//	// result 将包含: [0, 2, 4, 6, 8]
//
// 示例 2: 筛选长度大于5的字符串
//
//	words := xiter.FromSlice([]string{"apple", "banana", "cherry", "date"})
//	longWords := xiter.Filter(words, func(s string) bool { return len(s) > 5 })
//	result, _ := xiter.ToSlice(longWords)
//	// result 将包含: ["banana", "cherry"]
func Filter[E any](s iter.Seq[E], f func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if f(e) && !yield(e) {
				return
			}
		}
	}
}

// Filter2 返回一个仅包含满足谓词函数条件的键值对的新序列。
//
// 示例: 筛选值为偶数的键值对
//
//	input := xiter.Enumerate(xiter.Range(0, 5)) // 生成序列: (0,0), (1,1), (2,2), (3,3), (4,4)
//	evenValues := xiter.Filter2(input, func(k, v int) bool { return v%2 == 0 })
//	result, _ := xiter.ToSlice2(evenValues)
//	// result 将包含: [(0, 0), (2, 2), (4, 4)]
func Filter2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// Fold 通过对每个元素应用给定函数，将序列缩减为单个值。
//
// 示例 1: 计算序列元素总和
//
//	numbers := xiter.Range(1, 5)     // 生成序列: 1, 2, 3, 4
//	sum := xiter.Fold(numbers, 0, func(acc, e int) int { return acc + e })
//	// sum 将为: 10 (1+2+3+4)
//
// 示例 2: 连接字符串
//
//	words := xiter.FromSlice([]string{"Hello", " ", "World"})
//	phrase := xiter.Fold(words, "", func(acc, s string) string { return acc + s })
//	// phrase 将为: "Hello World"
func Fold[E any, A any](s iter.Seq[E], init A, f func(A, E) A) A {
	for e := range s {
		init = f(init, e)
	}
	return init
}

// Fold2 通过对每个键值对应用给定函数，将序列缩减为单个值。
//
// 示例: 计算所有值的总和
//
//	input := xiter.Enumerate(xiter.Range(1, 4)) // 生成序列: (0,1), (1,2), (2,3)
//	sum := xiter.Fold2(input, 0, func(acc int, k int, v int) int { return acc + v })
//	// sum 将为: 6 (1+2+3)
func Fold2[K, V, A any](s iter.Seq2[K, V], init A, f func(A, K, V) A) A {
	for k, v := range s {
		init = f(init, k, v)
	}
	return init
}

// Size 返回序列中的元素数量。
//
// 示例 1: 计算整数序列的长度
//
//	seq := xiter.Range(0, 5)       // 生成序列: 0, 1, 2, 3, 4
//	count := xiter.Size(seq)
//	// count 将为: 5
//
// 示例 2: 计算空序列的长度
//
//	emptySeq := xiter.FromSlice([]string{})
//	count := xiter.Size(emptySeq)
//	// count 将为: 0
func Size[E any](s iter.Seq[E]) int {
	return SizeFunc(s, func(_ E) bool { return true })
}

// Size2 返回键值对序列中的元素数量。
//
// 示例: 计算键值对序列的长度
//
//	input := xiter.Enumerate(xiter.Range(1, 4)) // 生成序列: (0,1), (1,2), (2,3)
//	count := xiter.Size2(input)
//	// count 将为: 3
func Size2[K, V any](s iter.Seq2[K, V]) int {
	return SizeFunc2(s, func(_ K, _ V) bool { return true })
}

// SizeFunc 返回满足谓词函数条件的元素数量。
//
// 示例 1: 计算偶数元素的数量
//
//	numbers := xiter.Range(0, 10)    // 生成序列: 0-9
//	evenCount := xiter.SizeFunc(numbers, func(e int) bool { return e%2 == 0 })
//	// evenCount 将为: 5
//
// 示例 2: 计算长度大于3的字符串数量
//
//	words := xiter.FromSlice([]string{"a", "bb", "ccc", "dddd"})
//	longCount := xiter.SizeFunc(words, func(s string) bool { return len(s) > 3 })
//	// longCount 将为: 1
func SizeFunc[E any](s iter.Seq[E], f func(E) bool) int {
	size := 0
	for e := range s {
		if f(e) {
			size++
		}
	}
	return size
}

// SizeFunc2 返回满足谓词函数条件的键值对数量。
//
// 示例: 计算值大于2的键值对数量
//
//	input := xiter.Enumerate(xiter.Range(1, 5)) // 生成序列: (0,1), (1,2), (2,3), (3,4)
//	count := xiter.SizeFunc2(input, func(k, v int) bool { return v > 2 })
//	// count 将为: 2 (键值对 (2,3) 和 (3,4))
func SizeFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) int {
	size := 0
	for k, v := range s {
		if f(k, v) {
			size++
		}
	}
	return size
}

// SizeValue 返回序列中与指定值相等的元素数量。
//
// 示例: 计算序列中值为5的元素数量
//
//	numbers := xiter.FromSlice([]int{2, 5, 5, 7, 5})
//	count := xiter.SizeValue(numbers, 5)
//	// count 将为: 3
func SizeValue[E comparable](s iter.Seq[E], v E) int {
	return SizeFunc(s, func(e E) bool { return e == v })
}

// SizeValue2 返回键值对序列中值等于指定值的元素数量。
//
// 示例: 计算值为2的键值对数量
//
//	pairs := xiter.Enumerate(xiter.Range(1, 4)) // 生成序列: (0,1), (1,2), (2,3)
//	count := xiter.SizeValue2(pairs, 2)
//	// count 将为: 1
func SizeValue2[K any, V comparable](s iter.Seq2[K, V], v V) int {
	return SizeFunc2(s, func(_ K, val V) bool { return val == v })
}

// Max 返回序列中的最大元素，使用 cmp.Compare 进行比较。如果序列为空，返回元素类型的零值和 false。
//
// 示例: 查找整数序列的最大值
//
//	numbers := xiter.FromSlice([]int{3, 1, 4, 1, 5, 9})
//	maxVal, ok := xiter.Max(numbers)
//	if ok {
//		fmt.Println("最大值:", maxVal) // 输出: 最大值: 9
//	}
func Max[E cmp.Ordered](s iter.Seq[E]) (E, bool) {
	return MaxFunc(s, cmp.Compare[E])
}

// MaxFunc 使用自定义比较函数返回序列中的最大元素。如果序列为空，返回元素类型的零值和 false。
//
// 示例: 查找字符串序列中最长的字符串
//
//	words := xiter.FromSlice([]string{"apple", "banana", "cherry"})
//	longest, ok := xiter.MaxFunc(words, func(a, b string) int {
//		return cmp.Compare(len(a), len(b))
//	})
//	if ok {
//		fmt.Println("最长字符串:", longest) // 输出: 最长字符串: banana
//	}
func MaxFunc[E any](s iter.Seq[E], cmp func(E, E) int) (E, bool) {
	it, stop := iter.Pull(s)
	defer stop()
	current, ok := it()
	if !ok {
		return current, false
	}
	for elem, ok := it(); ok; elem, ok = it() {
		if cmp(elem, current) > 0 {
			current = elem
		}
	}
	return current, true
}

// Min 返回序列中的最小元素，使用 cmp.Compare 进行比较。如果序列为空，返回元素类型的零值和 false。
//
// 示例: 查找整数序列的最小值
//
//	numbers := xiter.FromSlice([]int{5, 2, 7, 1, 3})
//	minVal, ok := xiter.Min(numbers)
//	if ok {
//		fmt.Println("最小值:", minVal) // 输出: 最小值: 1
//	}
func Min[E cmp.Ordered](s iter.Seq[E]) (E, bool) {
	return MinFunc(s, cmp.Compare[E])
}

// MinFunc 使用自定义比较函数返回序列中的最小元素。如果序列为空，返回元素类型的零值和 false。
//
// 示例: 查找字符串序列中最短的字符串
//
//	words := xiter.FromSlice([]string{"apple", "banana", "cherry", "date"})
//	shortest, ok := xiter.MinFunc(words, func(a, b string) int {
//		return cmp.Compare(len(a), len(b))
//	})
//	if ok {
//		fmt.Println("最短字符串:", shortest) // 输出: 最短字符串: date
//	}
func MinFunc[E any](s iter.Seq[E], cmp func(E, E) int) (E, bool) {
	it, stop := iter.Pull(s)
	defer stop()
	current, ok := it()
	if !ok {
		return current, false
	}
	for elem, ok := it(); ok; elem, ok = it() {
		if cmp(elem, current) < 0 {
			current = elem
		}
	}
	return current, true
}

// ForEach 对序列中的每个元素应用指定函数（无返回值，主要用于副作用操作如打印或修改外部状态）。
//
// 示例: 打印序列中的每个元素
//
//	seq := xiter.Range(1, 4)       // 生成序列: 1, 2, 3
//	xiter.ForEach(seq, func(e int) {
//		fmt.Printf("%d ", e)
//	})
//	// 输出: 1 2 3
func ForEach[E any](s iter.Seq[E], f func(E)) {
	for e := range s {
		f(e)
	}
}

// ForEach2 对键值对序列中的每个元素应用指定函数（无返回值，主要用于副作用操作）。
//
// 示例: 打印键值对
//
//	input := xiter.Enumerate(xiter.Range(1, 4)) // 生成序列: (0,1), (1,2), (2,3)
//	xiter.ForEach2(input, func(k, v int) {
//		fmt.Printf("(%d,%d) ", k, v)
//	})
//	// 输出: (0,1) (1,2) (2,3)
func ForEach2[K, V any](s iter.Seq2[K, V], f func(K, V)) {
	for k, v := range s {
		f(k, v)
	}
}

// FilterMap 对序列中的每个元素应用指定函数，该函数返回值和布尔值。仅当布尔值为true时，保留返回值到结果序列中。
//
// 示例: 筛选偶数并计算平方
//
//	numbers := xiter.Range(0, 5)    // 生成序列: 0, 1, 2, 3, 4
//	squares := xiter.FilterMap(numbers, func(e int) (int, bool) {
//		if e%2 == 0 {
//			return e*e, true
//		}
//		return 0, false
//	})
//	result, _ := xiter.ToSlice(squares)
//	// result 将包含: [0, 4, 16]
func FilterMap[E1, E2 any](s iter.Seq[E1], f func(E1) (E2, bool)) iter.Seq[E2] {
	return func(yield func(E2) bool) {
		for e1 := range s {
			if e2, ok := f(e1); ok && !yield(e2) {
				return
			}
		}
	}
}

// FilterMap2 对键值对序列中的每个元素应用指定函数，该函数返回键、值和布尔值。仅当布尔值为true时，保留键值对到结果序列中。
//
// 示例: 筛选值为偶数的键值对并将值加倍
//
//	input := xiter.Enumerate(xiter.Range(1, 5)) // 生成序列: (0,1), (1,2), (2,3), (3,4)
//	result := xiter.FilterMap2(input, func(k, v int) (int, int, bool) {
//		if v%2 == 0 {
//			return k, v*2, true
//		}
//		return 0, 0, false
//	})
//	// 结果序列将包含: (1,4), (3,8)
func FilterMap2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], f func(K1, V1) (K2, V2, bool)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range s {
			if k2, v2, ok := f(k1, v1); ok && !yield(k2, v2) {
				return
			}
		}
	}
}

// Keys 从键值对序列中提取所有键，返回一个包含这些键的序列。
//
// 示例: 提取键值对序列中的键
//
//	input := xiter.Enumerate(xiter.Range(1, 4)) // 生成序列: (0,1), (1,2), (2,3)
//	keys := xiter.Keys(input)
//	result, _ := xiter.ToSlice(keys)
//	// result 将包含: [0, 1, 2]
func Keys[K, V any](s iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range s {
			if !yield(k) {
				return
			}
		}
	}
}

// Values 从键值对序列中提取所有值，返回一个包含这些值的序列。
//
// 示例: 提取键值对序列中的值
//
//	input := xiter.Enumerate(xiter.Range(1, 4)) // 生成序列: (0,1), (1,2), (2,3)
//	values := xiter.Values(input)
//	result, _ := xiter.ToSlice(values)
//	// result 将包含: [1, 2, 3]
func Values[K, V any](s iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

// Concat 将多个序列连接成一个单一的序列。
//
// 示例: 连接两个整数序列
//
//	seq1 := xiter.Range(1, 3)       // 生成序列: 1, 2
//	seq2 := xiter.Range(4, 6)       // 生成序列: 4, 5
//	combined := xiter.Concat(seq1, seq2)
//	result, _ := xiter.ToSlice(combined)
//	// result 将包含: [1, 2, 4, 5]
func Concat[E any](seqs ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, s := range seqs {
			for e := range s {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Concat2 将多个键值对序列连接成一个单一的键值对序列。
//
// 示例: 连接两个键值对序列
//
//	seq1 := xiter.Enumerate(xiter.Range(1, 3)) // 生成序列: (0,1), (1,2)
//	seq2 := xiter.Enumerate(xiter.Range(4, 6)) // 生成序列: (0,4), (1,5)
//	combined := xiter.Concat2(seq1, seq2)
//	result, _ := xiter.ToSlice2(combined)
//	// result 将包含: [(0,1), (1,2), (0,4), (1,5)]
func Concat2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, s := range seqs {
			for k, v := range s {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// EqualFunc 使用指定的比较函数判断两个序列是否相等。
//
// 示例: 比较两个字符串序列是否忽略大小写相等
//
//	seq1 := xiter.FromSlice([]string{"Apple", "Banana"})
//	seq2 := xiter.FromSlice([]string{"apple", "banana"})
//	equal := xiter.EqualFunc(seq1, seq2, func(a, b string) bool {
//		return strings.ToLower(a) == strings.ToLower(b)
//	})
//	// equal 将为: true
func EqualFunc[E any](x, y iter.Seq[E], f func(E, E) bool) bool {
	it1, stop1 := iter.Pull(x)
	defer stop1()
	it2, stop2 := iter.Pull(y)
	defer stop2()
	for {
		e1, ok1 := it1()
		e2, ok2 := it2()
		if !ok1 && !ok2 {
			return true
		}
		if ok1 != ok2 || !f(e1, e2) {
			return false
		}
	}
}

// Equal 判断两个可比较类型的序列是否相等（元素数量和值均相同）。
//
// 示例: 比较两个整数序列是否相等
//
//	seq1 := xiter.FromSlice([]int{1, 2, 3})
//	seq2 := xiter.FromSlice([]int{1, 2, 3})
//	seq3 := xiter.FromSlice([]int{1, 3, 2})
//	fmt.Println(xiter.Equal(seq1, seq2)) // 输出: true
//	fmt.Println(xiter.Equal(seq1, seq3)) // 输出: false
func Equal[E comparable](x, y iter.Seq[E]) bool {
	return EqualFunc(x, y, func(e1, e2 E) bool { return e1 == e2 })
}

// EqualFunc2 使用指定的比较函数判断两个键值对序列是否相等。
//
// 示例: 比较两个键值对序列是否键相等且值忽略大小写相等
//
//	seq1 := xiter.FromSlice2([]iter.Pair[string, string]{{"a", "Apple"}, {"b", "Banana"}})
//	seq2 := xiter.FromSlice2([]iter.Pair[string, string]{{"a", "apple"}, {"b", "banana"}})
//	equal := xiter.EqualFunc2(seq1, seq2, func(k1, v1, k2, v2 string) bool {
//		return k1 == k2 && strings.ToLower(v1) == strings.ToLower(v2)
//	})
//	// equal 将为: true
func EqualFunc2[K, V any](x, y iter.Seq2[K, V], f func(K, V, K, V) bool) bool {
	it1, stop1 := iter.Pull2(x)
	defer stop1()
	it2, stop2 := iter.Pull2(y)
	defer stop2()
	for {
		k1, v1, ok1 := it1()
		k2, v2, ok2 := it2()
		if !ok1 && !ok2 {
			return true
		}
		if ok1 != ok2 || !f(k1, v1, k2, v2) {
			return false
		}
	}
}

// Equal2 判断两个键值对序列是否相等（键和值均为可比较类型且完全匹配）。
//
// 示例: 比较两个键值对序列是否相等
//
//	seq1 := xiter.FromSlice2([]iter.Pair[int, string]{{1, "a"}, {2, "b"}})
//	seq2 := xiter.FromSlice2([]iter.Pair[int, string]{{1, "a"}, {2, "b"}})
//	seq3 := xiter.FromSlice2([]iter.Pair[int, string]{{1, "a"}, {3, "b"}})
//	fmt.Println(xiter.Equal2(seq1, seq2)) // 输出: true
//	fmt.Println(xiter.Equal2(seq1, seq3)) // 输出: false
func Equal2[K, V comparable](x, y iter.Seq2[K, V]) bool {
	return EqualFunc2(x, y, func(k1 K, v1 V, k2 K, v2 V) bool { return k1 == k2 && v1 == v2 })
}

// Sum 计算序列中所有元素的总和（元素类型需支持加法运算）。
//
// 示例: 计算整数序列的总和
//
//	numbers := xiter.Range(1, 5) // 生成序列: 1, 2, 3, 4
//	total := xiter.Sum(numbers)
//	// total 将为: 10
func Sum[E numbric](s iter.Seq[E]) E {
	return Fold(s, E(0), func(a E, e E) E { return a + e })
}

// Join 将键值对序列转换为元素序列，使用指定函数将每个键值对映射为一个元素。
//
// 示例: 将键值对序列转换为字符串序列
//
//	input := xiter.FromSlice2([]iter.Pair[int, string]{{1, "a"}, {2, "b"}})
//	strs := xiter.Join(input, func(k int, v string) string { return fmt.Sprintf("%d=%s", k, v) })
//	result, _ := xiter.ToSlice(strs)
//	// result 将包含: ["1=a", "2=b"]
func Join[E, K, V any](s iter.Seq2[K, V], f func(K, V) E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Split 将元素序列转换为键值对序列，使用指定函数将每个元素拆分为键和值。
//
// 示例: 将字符串序列拆分为键值对序列
//
//	strs := xiter.FromSlice([]string{"1=a", "2=b"})
//	pairs := xiter.Split(strs, func(s string) (int, string) {
//		parts := strings.Split(s, "=")
//		k, _ := strconv.Atoi(parts[0])
//		return k, parts[1]
//	})
//	result, _ := xiter.ToSlice2(pairs)
//	// result 将包含: [(1, "a"), (2, "b")]
func Split[K, V, E any](s iter.Seq[E], f func(E) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := range s {
			k, v := f(e)
			if !yield(k, v) {
				return
			}
		}
	}
}

// Any 判断序列中是否存在满足谓词函数的元素。
//
// 示例: 检查序列中是否存在偶数
//
//	numbers := xiter.Range(1, 5)     // 生成序列: 1, 2, 3, 4
//	hasEven := xiter.Any(numbers, func(e int) bool { return e%2 == 0 })
//	// hasEven 将为: true
func Any[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

// Any2 判断键值对序列中是否存在满足谓词函数的元素。
//
// 示例: 检查是否存在值为偶数的键值对
//
//	input := xiter.Enumerate(xiter.Range(1, 5)) // 生成序列: (0,1), (1,2), (2,3), (3,4)
//	hasEven := xiter.Any2(input, func(k, v int) bool { return v%2 == 0 })
//	// hasEven 将为: true
func Any2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return true
		}
	}
	return false
}

// All 判断序列中所有元素是否都满足谓词函数。
//
// 示例: 检查序列中所有元素是否都是正数
//
//	numbers := xiter.FromSlice([]int{2, 4, 6, 8})
//	allPositive := xiter.All(numbers, func(e int) bool { return e > 0 })
//	// allPositive 将为: true
func All[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if !f(e) {
			return false
		}
	}
	return true
}

// All2 判断键值对序列中所有元素是否都满足谓词函数。
//
// 示例: 检查所有值是否都大于0
//
//	input := xiter.Enumerate(xiter.Range(1, 4)) // 生成序列: (0,1), (1,2), (2,3)
//	allPositive := xiter.All2(input, func(k, v int) bool { return v > 0 })
//	// allPositive 将为: true
func All2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if !f(k, v) {
			return false
		}
	}
	return true
}

// Cast 将 any 类型的序列转换为指定类型 E 的序列，并返回转换结果和成功标志。
//
// 示例: 将 any 序列转换为 int 序列
//
//	values := xiter.FromSlice([]any{1, "two", 3, "four"})
//	casted := xiter.Cast[int](values)
//	result, _ := xiter.ToSlice2(casted)
//	// result 将包含: [(1, true), (0, false), (3, true), (0, false)]
func Cast[E any](s iter.Seq[any]) iter.Seq2[E, bool] {
	return func(yield func(E, bool) bool) {
		for a := range s {
			if e, ok := a.(E); !yield(e, ok) {
				return
			}
		}
	}
}

// Distinct 返回序列中所有不重复的元素，保持首次出现的顺序。
//
// 示例: 移除序列中的重复元素
//
//	numbers := xiter.FromSlice([]int{1, 2, 2, 3, 3, 3})
//	unique := xiter.Distinct(numbers)
//	result, _ := xiter.ToSlice(unique)
//	// result 将包含: [1, 2, 3]
func Distinct[E comparable](s iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		m := make(map[E]struct{})
		for e := range s {
			if _, ok := m[e]; !ok {
				m[e] = struct{}{}
				if !yield(e) {
					return
				}
			}
		}
	}
}

// ContainsFunc 判断序列中是否包含满足谓词函数的元素。
//
// 示例: 检查序列中是否包含长度大于3的字符串
//
//	words := xiter.FromSlice([]string{"a", "bb", "ccc", "dddd"})
//	hasLong := xiter.ContainsFunc(words, func(s string) bool { return len(s) > 3 })
//	// hasLong 将为: true
func ContainsFunc[E any](s iter.Seq[E], f func(E) bool) bool {
	for e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

// Contains 判断序列中是否包含指定的可比较元素。
//
// 示例: 检查序列中是否包含元素 3
//
//	numbers := xiter.FromSlice([]int{1, 2, 3, 4})
//	hasThree := xiter.Contains(numbers, 3)
//	// hasThree 将为: true
func Contains[E comparable](s iter.Seq[E], v E) bool {
	return ContainsFunc(s, func(e E) bool { return e == v })
}

func ContainsFunc2[K, V any](s iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range s {
		if f(k, v) {
			return true
		}
	}
	return false
}

// Contains2 判断键值对序列中是否包含指定的键和值。
//
// 示例: 检查是否包含键为1且值为"a"的键值对
//
//	pairs := xiter.FromSlice2([]iter.Pair[int, string]{{0, "x"}, {1, "a"}, {2, "b"}})
//	contains := xiter.Contains2(pairs, 1, "a")
//	// contains 将为: true
func Contains2[K, V comparable](s iter.Seq2[K, V], k K, v V) bool {
	return ContainsFunc2(s, func(ck K, cv V) bool { return ck == k && cv == v })
}

// IsSortedFunc 使用指定的比较函数判断序列是否按统一顺序排列（升序或降序）。
//
// 示例: 判断序列是否按降序排列
//
//	numbers := xiter.FromSlice([]int{5, 4, 3, 2, 1})
//	isDesc := xiter.IsSortedFunc(numbers, func(a, b int) int { return cmp.Compare(b, a) })
//	// isDesc 将为: true
func IsSortedFunc[E any](s iter.Seq[E], f func(E, E) int) bool {
	it, stop := iter.Pull(s)
	defer stop()
	first, fOk := it()
	second, sOk := it()
	if !fOk || !sOk {
		return true
	}
	prev := second
	initOrder := f(first, second)
	for x, ok := it(); ok; x, ok = it() {
		order := f(prev, x)
		if initOrder == 0 {
			initOrder = order
		}
		if order != initOrder && order != 0 {
			return false
		}
		prev = x
	}
	return true
}

// IsSorted 判断有序类型的序列是否按升序排列。
//
// 示例: 检查整数序列是否升序排列
//
//	seq1 := xiter.FromSlice([]int{1, 2, 3, 4})
//	seq2 := xiter.FromSlice([]int{1, 3, 2, 4})
//	fmt.Println(xiter.IsSorted(seq1)) // 输出: true
//	fmt.Println(xiter.IsSorted(seq2)) // 输出: false
func IsSorted[E cmp.Ordered](s iter.Seq[E]) bool {
	return IsSortedFunc(s, cmp.Compare[E])
}

// Zip 将两个元素序列配对为一个键值对序列，迭代至任一序列耗尽为止。
func Zip[E1, E2 any](a iter.Seq[E1], b iter.Seq[E2]) iter.Seq2[E1, E2] {
    return func(yield func(E1, E2) bool) {
        ia, stopA := iter.Pull(a)
        defer stopA()
        ib, stopB := iter.Pull(b)
        defer stopB()
        for {
            va, oka := ia()
            vb, okb := ib()
            if !oka || !okb {
                return
            }
            if !yield(va, vb) {
                return
            }
        }
    }
}

// ZipWith 使用给定函数将两个序列逐项合并为一个新序列，迭代至任一序列耗尽为止。
func ZipWith[E1, E2, R any](a iter.Seq[E1], b iter.Seq[E2], f func(E1, E2) R) iter.Seq[R] {
    return func(yield func(R) bool) {
        ia, stopA := iter.Pull(a)
        defer stopA()
        ib, stopB := iter.Pull(b)
        defer stopB()
        for {
            va, oka := ia()
            vb, okb := ib()
            if !oka || !okb {
                return
            }
            if !yield(f(va, vb)) {
                return
            }
        }
    }
}

// ZipWith2 使用给定函数将两个键值对序列逐项合并为一个新的键值对序列，迭代至任一序列耗尽为止。
func ZipWith2[K1, V1, K2, V2, K3, V3 any](a iter.Seq2[K1, V1], b iter.Seq2[K2, V2], f func(K1, V1, K2, V2) (K3, V3)) iter.Seq2[K3, V3] {
    return func(yield func(K3, V3) bool) {
        ia, stopA := iter.Pull2(a)
        defer stopA()
        ib, stopB := iter.Pull2(b)
        defer stopB()
        for {
            ka, va, oka := ia()
            kb, vb, okb := ib()
            if !oka || !okb {
                return
            }
            k3, v3 := f(ka, va, kb, vb)
            if !yield(k3, v3) {
                return
            }
        }
    }
}
