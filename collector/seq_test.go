package collector

import (
	"iter"
	"reflect"
	"slices"
	"testing"

	"github.com/go-board/xiter"
)

func TestCollect(t *testing.T) {
	got := Collect(seqOf(1, 2, 3), ToSlice[int]())
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v, want [1 2 3]", got)
	}
}

func TestToSlice(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		got := ToSlice[int]()(seqOf(1, 2, 3))
		if !slices.Equal(got, []int{1, 2, 3}) {
			t.Fatalf("got %v", got)
		}
	})
	t.Run("empty", func(t *testing.T) {
		got := ToSlice[int]()(seqOf[int]())
		if got != nil {
			t.Fatalf("got %v, want nil", got)
		}
	})
}

func TestToSet(t *testing.T) {
	got := ToSet[int]()(seqOf(1, 2, 2, 3, 1))
	want := map[int]struct{}{1: {}, 2: {}, 3: {}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got := ToSet[int]()(seqOf[int]()); len(got) != 0 {
		t.Fatalf("got %v, want empty", got)
	}
}

func TestToMap(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		got := ToMap(func(e int) (int, int) {
			return e, e * e
		})(seqOf(1, 2, 3))
		want := map[int]int{1: 1, 2: 4, 3: 9}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("duplicate last wins", func(t *testing.T) {
		got := ToMap(func(e int) (int, int) {
			return e % 2, e
		})(seqOf(1, 3, 2, 4))
		want := map[int]int{1: 3, 0: 4}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("empty", func(t *testing.T) {
		got := ToMap(func(e int) (int, int) { return e, e })(seqOf[int]())
		if len(got) != 0 {
			t.Fatalf("got %v, want empty", got)
		}
	})
}

func TestToMapMerge(t *testing.T) {
	got := ToMapMerge(
		func(e int) (int, int) { return e % 2, e },
		func(a, b int) int { return a + b },
	)(seqOf(1, 3, 2, 4))
	want := map[int]int{1: 4, 0: 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestJoining(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		if got := Joining(", ")(seqOf("a", "b", "c")); got != "a, b, c" {
			t.Fatalf("got %q", got)
		}
	})
	t.Run("empty", func(t *testing.T) {
		if got := Joining(", ")(seqOf[string]()); got != "" {
			t.Fatalf("got %q, want empty", got)
		}
	})
	t.Run("single", func(t *testing.T) {
		if got := Joining(", ")(seqOf("a")); got != "a" {
			t.Fatalf("got %q", got)
		}
	})
}

func TestGroupingBy(t *testing.T) {
	got := GroupingBy(func(e int) int { return e % 2 })(seqOf(1, 2, 3, 4, 5))
	want := map[int][]int{0: {2, 4}, 1: {1, 3, 5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got := GroupingBy(func(e int) int { return e % 2 })(seqOf[int]()); len(got) != 0 {
		t.Fatalf("got %v, want empty", got)
	}
}

// counting wraps xiter.Size as a Collector, demonstrating how xiter terminal
// operations can be reused downstream. This mirrors the pattern documented on
// GroupingByDownstream.
func counting[E any]() Collector[E, int] {
	return Collector[E, int](func(s iter.Seq[E]) int { return xiter.Size(s) })
}

func TestGroupingByDownstream(t *testing.T) {
	t.Run("with counting via xiter.Size", func(t *testing.T) {
		got := GroupingByDownstream(
			func(e int) int { return e % 2 },
			counting[int](),
		)(seqOf(1, 2, 3, 4, 5))
		want := map[int]int{0: 2, 1: 3}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("with toslice", func(t *testing.T) {
		got := GroupingByDownstream(
			func(e int) int { return e % 2 },
			ToSlice[int](),
		)(seqOf(1, 2, 3, 4))
		want := map[int][]int{0: {2, 4}, 1: {1, 3}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("empty", func(t *testing.T) {
		got := GroupingByDownstream(
			func(e int) int { return e % 2 },
			counting[int](),
		)(seqOf[int]())
		if len(got) != 0 {
			t.Fatalf("got %v, want empty", got)
		}
	})
	t.Run("downstream early stop exercises fromSlice yield-false branch", func(t *testing.T) {
		// A downstream collector that breaks after the first element forces
		// fromSlice's yield to return false on the next call, covering the
		// !yield early-termination guard.
		takeFirst := Collector[int, []int](func(s iter.Seq[int]) []int {
			var out []int
			for e := range s {
				out = append(out, e)
				break
			}
			return out
		})
		got := GroupingByDownstream(
			func(e int) int { return e % 2 },
			takeFirst,
		)(seqOf(1, 3, 5, 2, 4))
		want := map[int][]int{1: {1}, 0: {2}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}

func TestPartitioningBy(t *testing.T) {
	t.Run("both sides", func(t *testing.T) {
		got := PartitioningBy(func(e int) bool { return e%2 == 0 })(seqOf(1, 2, 3, 4, 5))
		if !slices.Equal(got.Pass, []int{2, 4}) || !slices.Equal(got.Fail, []int{1, 3, 5}) {
			t.Fatalf("Pass=%v Fail=%v", got.Pass, got.Fail)
		}
	})
	t.Run("all pass", func(t *testing.T) {
		got := PartitioningBy(func(e int) bool { return true })(seqOf(1, 2, 3))
		if !slices.Equal(got.Pass, []int{1, 2, 3}) || got.Fail != nil {
			t.Fatalf("Pass=%v Fail=%v", got.Pass, got.Fail)
		}
	})
	t.Run("all fail", func(t *testing.T) {
		got := PartitioningBy(func(e int) bool { return false })(seqOf(1, 2, 3))
		if got.Pass != nil || !slices.Equal(got.Fail, []int{1, 2, 3}) {
			t.Fatalf("Pass=%v Fail=%v", got.Pass, got.Fail)
		}
	})
	t.Run("empty", func(t *testing.T) {
		got := PartitioningBy(func(e int) bool { return true })(seqOf[int]())
		if got.Pass != nil || got.Fail != nil {
			t.Fatalf("Pass=%v Fail=%v", got.Pass, got.Fail)
		}
	})
}
