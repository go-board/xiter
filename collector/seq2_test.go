package collector

import (
	"reflect"
	"slices"
	"testing"
)

func TestCollect2(t *testing.T) {
	got := Collect2(seq2From([]int{1, 2}, []string{"a", "b"}), ToMap2[int, string]())
	want := map[int]string{1: "a", 2: "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestToMap2(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		got := ToMap2[int, string]()(seq2From([]int{1, 2}, []string{"a", "b"}))
		want := map[int]string{1: "a", 2: "b"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("duplicate last wins", func(t *testing.T) {
		got := ToMap2[int, string]()(seq2From([]int{1, 1, 2}, []string{"a", "b", "c"}))
		want := map[int]string{1: "b", 2: "c"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("empty", func(t *testing.T) {
		got := ToMap2[int, string]()(seq2From[int, string](nil, nil))
		if len(got) != 0 {
			t.Fatalf("got %v, want empty", got)
		}
	})
}

func TestToMap2Merge(t *testing.T) {
	got := ToMap2Merge[int, int](func(a, b int) int { return a + b })(
		seq2From([]int{1, 1, 2}, []int{10, 20, 30}),
	)
	want := map[int]int{1: 30, 2: 30}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestToKeys(t *testing.T) {
	got := ToKeys[int, string]()(seq2From([]int{1, 2, 3}, []string{"a", "b", "c"}))
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
	if got := ToKeys[int, string]()(seq2From[int, string](nil, nil)); got != nil {
		t.Fatalf("got %v, want nil", got)
	}
}

func TestToValues(t *testing.T) {
	got := ToValues[int, string]()(seq2From([]int{1, 2, 3}, []string{"a", "b", "c"}))
	if !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Fatalf("got %v", got)
	}
	if got := ToValues[int, string]()(seq2From[int, string](nil, nil)); got != nil {
		t.Fatalf("got %v, want nil", got)
	}
}
