package stream

import (
	"cmp"
	"iter"
	"reflect"
	"testing"

	"github.com/go-board/xiter"
)

func collect[E any](s Stream[E]) []E {
	out := make([]E, 0)
	s.ForEach(func(v E) { out = append(out, v) })
	return out
}

func collect2[K, V any](s Stream2[K, V]) []struct {
	K K
	V V
} {
	out := make([]struct {
		K K
		V V
	}, 0)
	s.ForEach(func(k K, v V) {
		out = append(out, struct {
			K K
			V V
		}{k, v})
	})
	return out
}

func TestStreamMethodsCoverage(t *testing.T) {
	base := Of(xiter.Range1(6))
	if got := collect(Of(base.Seq())); !reflect.DeepEqual(got, []int{0, 1, 2, 3, 4, 5}) {
		t.Fatalf("Seq/Of failed: %v", got)
	}

	if got := collect(base.Filter(func(v int) bool { return v%2 == 0 })); !reflect.DeepEqual(got, []int{0, 2, 4}) {
		t.Fatalf("Filter: %v", got)
	}
	if got := collect(base.Map(func(v int) int { return v + 1 })); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5, 6}) {
		t.Fatalf("Map: %v", got)
	}
	if got := collect(base.MapWhile(func(v int) (int, bool) { return v * 2, v < 3 })); !reflect.DeepEqual(got, []int{0, 2, 4}) {
		t.Fatalf("MapWhile: %v", got)
	}
	if got := collect(base.FilterMap(func(v int) (int, bool) { return v * 10, v%2 == 1 })); !reflect.DeepEqual(got, []int{10, 30, 50}) {
		t.Fatalf("FilterMap: %v", got)
	}
	if got := collect(base.FlatMap(func(v int) iter.Seq[int] { return xiter.Once(v) }).Take(3)); !reflect.DeepEqual(got, []int{0, 1, 2}) {
		t.Fatalf("FlatMap: %v", got)
	}
	if got := collect(base.Take(2)); !reflect.DeepEqual(got, []int{0, 1}) {
		t.Fatalf("Take: %v", got)
	}
	if got := collect(base.Skip(4)); !reflect.DeepEqual(got, []int{4, 5}) {
		t.Fatalf("Skip: %v", got)
	}
	if got := collect(base.TakeWhile(func(v int) bool { return v < 3 })); !reflect.DeepEqual(got, []int{0, 1, 2}) {
		t.Fatalf("TakeWhile: %v", got)
	}
	if got := collect(base.SkipWhile(func(v int) bool { return v < 3 })); !reflect.DeepEqual(got, []int{3, 4, 5}) {
		t.Fatalf("SkipWhile: %v", got)
	}
	if got := collect(base.Chain(Of(xiter.Range2(6, 8)))); !reflect.DeepEqual(got, []int{0, 1, 2, 3, 4, 5, 6, 7}) {
		t.Fatalf("Chain: %v", got)
	}
	if got := collect2(base.Enumerate().Take(2)); !reflect.DeepEqual(got, []struct {
		K int
		V int
	}{{0, 0}, {1, 1}}) {
		t.Fatalf("Enumerate: %v", got)
	}

	sum := base.Fold(0, func(a, b int) int { return a + b })
	if sum != 15 || base.Size() != 6 || base.SizeFunc(func(v int) bool { return v%2 == 0 }) != 3 {
		t.Fatalf("Fold/Size: %d %d", sum, base.Size())
	}
	if !base.Any(func(v int) bool { return v == 5 }) || !base.All(func(v int) bool { return v >= 0 }) {
		t.Fatal("Any/All failed")
	}
	if f, ok := base.First(); !ok || f != 0 {
		t.Fatalf("First: %v %v", f, ok)
	}
	if f, ok := base.FirstFunc(func(v int) bool { return v > 2 }); !ok || f != 3 {
		t.Fatalf("FirstFunc: %v %v", f, ok)
	}
	if l, ok := base.Last(); !ok || l != 5 {
		t.Fatalf("Last: %v %v", l, ok)
	}
	if l, ok := base.LastFunc(func(v int) bool { return v%2 == 0 }); !ok || l != 4 {
		t.Fatalf("LastFunc: %v %v", l, ok)
	}
	if p, ok := base.Position(func(v int) bool { return v == 4 }); !ok || p != 4 {
		t.Fatalf("Position: %v %v", p, ok)
	}
	if !base.IsSortedFunc(cmp.Compare[int]) {
		t.Fatal("IsSortedFunc failed")
	}
	if c := base.CompareFunc(Of(xiter.Range1(6)), cmp.Compare[int]); c != 0 {
		t.Fatalf("CompareFunc: %d", c)
	}
	if !base.EqualFunc(Of(xiter.Range1(6)), func(a, b int) bool { return a == b }) {
		t.Fatal("EqualFunc failed")
	}
	if m, ok := base.MaxFunc(cmp.Compare[int]); !ok || m != 5 {
		t.Fatalf("MaxFunc: %v %v", m, ok)
	}
	if m, ok := base.MinFunc(cmp.Compare[int]); !ok || m != 0 {
		t.Fatalf("MinFunc: %v %v", m, ok)
	}
	if !base.ContainsFunc(func(v int) bool { return v == 2 }) {
		t.Fatal("ContainsFunc failed")
	}
}

func TestStream2MethodsCoverage(t *testing.T) {
	s1 := Of2(func(yield func(string, int) bool) {
		for _, kv := range []struct {
			k string
			v int
		}{{"a", 1}, {"b", 2}, {"c", 3}} {
			if !yield(kv.k, kv.v) {
				return
			}
		}
	})
	if got := collect2(Of2(s1.Seq2())); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"a", 1}, {"b", 2}, {"c", 3}}) {
		t.Fatalf("Seq2/Of2 failed: %v", got)
	}
	if got := collect2(s1.Filter(func(_ string, v int) bool { return v >= 2 })); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"b", 2}, {"c", 3}}) {
		t.Fatalf("Filter2: %v", got)
	}
	if got := collect2(s1.Map(func(k string, v int) (string, int) { return k + k, v * 10 })); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"aa", 10}, {"bb", 20}, {"cc", 30}}) {
		t.Fatalf("Map2: %v", got)
	}
	if got := collect2(s1.FilterMap(func(k string, v int) (string, int, bool) { return k, v + 1, v > 1 })); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"b", 3}, {"c", 4}}) {
		t.Fatalf("FilterMap2: %v", got)
	}
	if got := collect2(s1.MapWhile(func(k string, v int) (string, int, bool) { return k, v, v < 3 })); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"a", 1}, {"b", 2}}) {
		t.Fatalf("MapWhile2: %v", got)
	}
	if got := collect2(s1.Take(2)); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"a", 1}, {"b", 2}}) {
		t.Fatalf("Take2: %v", got)
	}
	if got := collect2(s1.Skip(1)); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"b", 2}, {"c", 3}}) {
		t.Fatalf("Skip2: %v", got)
	}
	if got := collect2(s1.TakeWhile(func(_ string, v int) bool { return v < 3 })); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"a", 1}, {"b", 2}}) {
		t.Fatalf("TakeWhile2: %v", got)
	}
	if got := collect2(s1.SkipWhile(func(_ string, v int) bool { return v < 2 })); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"b", 2}, {"c", 3}}) {
		t.Fatalf("SkipWhile2: %v", got)
	}
	s2 := Of2(func(yield func(string, int) bool) { _ = yield("d", 4) })
	if got := collect2(s1.Chain(s2)); !reflect.DeepEqual(got, []struct {
		K string
		V int
	}{{"a", 1}, {"b", 2}, {"c", 3}, {"d", 4}}) {
		t.Fatalf("Chain2: %v", got)
	}

	if s1.Size() != 3 || s1.SizeFunc(func(_ string, v int) bool { return v >= 2 }) != 2 || !s1.Any(func(_ string, v int) bool { return v == 2 }) || !s1.All(func(_ string, v int) bool { return v > 0 }) {
		t.Fatal("Size/Any/All2 failed")
	}
	if k, v, ok := s1.First(); !ok || k != "a" || v != 1 {
		t.Fatalf("First2: %v %v %v", k, v, ok)
	}
	if k, v, ok := s1.FirstFunc(func(k string, _ int) bool { return k == "b" }); !ok || k != "b" || v != 2 {
		t.Fatalf("FirstFunc2: %v %v %v", k, v, ok)
	}
	if k, v, ok := s1.Last(); !ok || k != "c" || v != 3 {
		t.Fatalf("Last2: %v %v %v", k, v, ok)
	}
	if k, v, ok := s1.LastFunc(func(_ string, v int) bool { return v >= 2 }); !ok || k != "c" || v != 3 {
		t.Fatalf("LastFunc2: %v %v %v", k, v, ok)
	}
	if p, ok := s1.Position(func(k string, _ int) bool { return k == "b" }); !ok || p != 1 {
		t.Fatalf("Position2: %v %v", p, ok)
	}
	if c := s1.CompareFunc(Of2(s1.Seq2()), func(k1 string, v1 int, k2 string, v2 int) int {
		if ck := cmp.Compare(k1, k2); ck != 0 {
			return ck
		}
		return cmp.Compare(v1, v2)
	}); c != 0 {
		t.Fatalf("CompareFunc2: %v", c)
	}
	if !s1.EqualFunc(Of2(s1.Seq2()), func(k1 string, v1 int, k2 string, v2 int) bool { return k1 == k2 && v1 == v2 }) {
		t.Fatal("EqualFunc2 failed")
	}
	if !s1.ContainsFunc(func(k string, v int) bool { return k == "c" && v == 3 }) {
		t.Fatal("ContainsFunc2 failed")
	}
}
