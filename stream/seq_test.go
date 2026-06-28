//go:build go1.27

package stream

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/go-board/xiter"
)

func TestSeqChainWithGenericMethods(t *testing.T) {
	got := Of(xiter.Range1(6)).
		Skip(1).
		Take(4).
		Filter(func(v int) bool { return v%2 == 0 }).
		Map(func(v int) string { return fmt.Sprintf("n=%d", v) }).
		Fold([]string{}, func(acc []string, v string) []string {
			return append(acc, v)
		})

	want := []string{"n=2", "n=4"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSeqFromFuncAndReduce(t *testing.T) {
	i := 0
	s := FromFunc(func() (int, bool) {
		if i >= 5 {
			return 0, false
		}
		i++
		return i, true
	})

	got, ok := s.Reduce(func(acc, v int) int {
		return acc + v
	})
	if !ok {
		t.Fatal("Reduce returned ok=false, want true")
	}
	if got != 15 {
		t.Fatalf("got %d, want 15", got)
	}

	empty, ok := Of(xiter.Empty[int]()).Reduce(func(acc, v int) int {
		return acc + v
	})
	if ok {
		t.Fatal("empty Reduce returned ok=true, want false")
	}
	if empty != 0 {
		t.Fatalf("got %d, want zero", empty)
	}

	got, ok, err := Of(xiter.Range1(5)).TryReduce(func(acc, v int) (int, error) {
		acc += v
		if v == 3 {
			return acc, errors.New("stop")
		}
		return acc, nil
	})
	if err == nil {
		t.Fatal("TryReduce returned nil error, want non-nil")
	}
	if !ok {
		t.Fatal("TryReduce returned ok=false, want true")
	}
	if got != 6 {
		t.Fatalf("got %d, want 6", got)
	}
}

func TestSeqTryMethods(t *testing.T) {
	wantErr := errors.New("stop")
	visited := 0

	err := Of(xiter.Range1(5)).TryForEach(func(v int) error {
		visited++
		if v == 2 {
			return wantErr
		}
		return nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if visited != 3 {
		t.Fatalf("visited %d elements, want 3", visited)
	}

	got, err := Of(xiter.Range1(5)).TryFold(0, func(acc, v int) (int, error) {
		acc += v
		if v == 3 {
			return acc, wantErr
		}
		return acc, nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if got != 6 {
		t.Fatalf("got accumulator %d, want 6", got)
	}
}

func TestSeqInspect(t *testing.T) {
	seen := []int{}
	got := Of(xiter.Range1(5)).
		Inspect(func(v int) { seen = append(seen, v) }).
		Take(2).
		Fold([]int{}, func(acc []int, v int) []int {
			return append(acc, v)
		})
	want := []int{0, 1}

	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if !slices.Equal(seen, want) {
		t.Fatalf("seen %v, want %v", seen, want)
	}
}

func TestSeqZipMethods(t *testing.T) {
	got := Of(xiter.Range1(5)).
		Zip(Of(xiter.Range2(10, 13))).
		Join(func(a, b int) string {
			return fmt.Sprintf("%d:%d", a, b)
		}).
		Fold([]string{}, func(acc []string, v string) []string {
			return append(acc, v)
		})
	want := []string{"0:10", "1:11", "2:12"}

	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	with := Of(xiter.Range1(5)).
		ZipWith(Of(xiter.Range2(10, 13)), func(a, b int) int {
			return a + b
		}).
		Fold([]int{}, func(acc []int, v int) []int {
			return append(acc, v)
		})
	wantWith := []int{10, 12, 14}

	if !slices.Equal(with, wantWith) {
		t.Fatalf("got %v, want %v", with, wantWith)
	}
}

func TestSeq2ChainWithGenericMethods(t *testing.T) {
	got := Of(xiter.Range1(5)).
		Split(func(v int) (string, int) {
			return fmt.Sprintf("k%d", v), v * 10
		}).
		Filter(func(_ string, v int) bool { return v >= 20 }).
		Map(func(k string, v int) (int, string) {
			return v / 10, k
		}).
		Join(func(k int, v string) string {
			return fmt.Sprintf("%d:%s", k, v)
		}).
		Fold([]string{}, func(acc []string, v string) []string {
			return append(acc, v)
		})

	want := []string{"2:k2", "3:k3", "4:k4"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSeq2FromFuncAndReduce(t *testing.T) {
	i := 0
	s := FromFunc2(func() (int, string, bool) {
		if i >= 3 {
			return 0, "", false
		}
		i++
		return i, fmt.Sprintf("v%d", i), true
	})

	gotK, gotV, ok := s.Reduce(func(accK int, accV string, k int, v string) (int, string) {
		return accK + k, accV + v
	})
	if !ok {
		t.Fatal("Reduce returned ok=false, want true")
	}
	if gotK != 6 || gotV != "v1v2v3" {
		t.Fatalf("got (%d, %q), want (6, %q)", gotK, gotV, "v1v2v3")
	}

	emptyK, emptyV, ok := Of2(xiter.Empty2[int, string]()).Reduce(func(accK int, accV string, k int, v string) (int, string) {
		return accK + k, accV + v
	})
	if ok {
		t.Fatal("empty Reduce returned ok=true, want false")
	}
	if emptyK != 0 || emptyV != "" {
		t.Fatalf("got (%d, %q), want zero values", emptyK, emptyV)
	}

	tryK, tryV, ok, err := Of(xiter.Range1(5)).Enumerate().TryReduce(func(accK, accV, k, v int) (int, int, error) {
		accK += k
		accV += v
		if k == 3 {
			return accK, accV, errors.New("stop")
		}
		return accK, accV, nil
	})
	if err == nil {
		t.Fatal("TryReduce returned nil error, want non-nil")
	}
	if !ok {
		t.Fatal("TryReduce returned ok=false, want true")
	}
	if tryK != 6 || tryV != 6 {
		t.Fatalf("got (%d, %d), want (6, 6)", tryK, tryV)
	}
}

func TestSeq2TryMethods(t *testing.T) {
	wantErr := errors.New("stop")
	visited := 0
	pairs := Of(xiter.Range1(5)).Enumerate()

	err := pairs.TryForEach(func(k, v int) error {
		visited++
		if k == 2 || v == 2 {
			return wantErr
		}
		return nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if visited != 3 {
		t.Fatalf("visited %d pairs, want 3", visited)
	}

	got, err := Of(xiter.Range1(5)).Enumerate().TryFold(0, func(acc, k, v int) (int, error) {
		acc += k + v
		if k == 3 {
			return acc, wantErr
		}
		return acc, nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if got != 12 {
		t.Fatalf("got accumulator %d, want 12", got)
	}
}

func TestSeq2Inspect(t *testing.T) {
	seen := map[int]int{}
	got := Of(xiter.Range1(5)).
		Enumerate().
		Inspect(func(k, v int) { seen[k] = v }).
		Take(2).
		Join(func(k, v int) int { return k + v }).
		Fold([]int{}, func(acc []int, v int) []int {
			return append(acc, v)
		})
	want := []int{0, 2}

	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if len(seen) != 2 || seen[0] != 0 || seen[1] != 1 {
		t.Fatalf("seen %v, want map[0:0 1:1]", seen)
	}
}

func TestSeq2KeyValueMethods(t *testing.T) {
	pairs := Of(xiter.Range2(10, 13)).Enumerate()

	keys := pairs.Keys().Fold([]int{}, func(acc []int, v int) []int {
		return append(acc, v)
	})
	wantKeys := []int{0, 1, 2}
	if !slices.Equal(keys, wantKeys) {
		t.Fatalf("keys %v, want %v", keys, wantKeys)
	}

	values := Of(xiter.Range2(10, 13)).Enumerate().Values().Fold([]int{}, func(acc []int, v int) []int {
		return append(acc, v)
	})
	wantValues := []int{10, 11, 12}
	if !slices.Equal(values, wantValues) {
		t.Fatalf("values %v, want %v", values, wantValues)
	}

	swapped := Of(xiter.Range2(10, 13)).
		Enumerate().
		Swap().
		Join(func(k, v int) string {
			return fmt.Sprintf("%d:%d", k, v)
		}).
		Fold([]string{}, func(acc []string, v string) []string {
			return append(acc, v)
		})
	wantSwapped := []string{"10:0", "11:1", "12:2"}
	if !slices.Equal(swapped, wantSwapped) {
		t.Fatalf("swapped %v, want %v", swapped, wantSwapped)
	}
}
