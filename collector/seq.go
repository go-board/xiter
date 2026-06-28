package collector

import (
	"iter"
	"strings"
)

// Partition holds the two halves produced by PartitioningBy. Either field may
// be nil when no element landed on that side.
type Partition[E any] struct {
	Pass []E // elements for which the predicate returned true
	Fail []E // elements for which the predicate returned false
}

// fromSlice turns a slice into an iter.Seq for internal use, e.g. feeding
// grouped elements into a downstream collector.
func fromSlice[E any](s []E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, e := range s {
			if !yield(e) {
				return
			}
		}
	}
}

// ToSlice returns a collector that accumulates elements into a slice in
// iteration order. An empty input yields a nil slice.
func ToSlice[E any]() Collector[E, []E] {
	return func(s iter.Seq[E]) []E {
		var out []E
		for e := range s {
			out = append(out, e)
		}
		return out
	}
}

// ToSet returns a collector that accumulates elements into a set represented as
// a map[E]struct{}. Duplicate elements collapse; only presence is retained. An
// empty input yields an empty, non-nil map.
func ToSet[E comparable]() Collector[E, map[E]struct{}] {
	return func(s iter.Seq[E]) map[E]struct{} {
		out := make(map[E]struct{})
		for e := range s {
			out[e] = struct{}{}
		}
		return out
	}
}

// ToMap returns a collector that builds a map by applying f to each element to
// produce its (key, value) pair. On duplicate keys the later element overwrites
// the earlier value (last wins); for custom duplicate handling use ToMapMerge.
// An empty input yields an empty, non-nil map.
func ToMap[E any, K comparable, V any](f func(E) (K, V)) Collector[E, map[K]V] {
	return func(s iter.Seq[E]) map[K]V {
		out := make(map[K]V)
		for e := range s {
			k, v := f(e)
			out[k] = v
		}
		return out
	}
}

// ToMapMerge is like ToMap but resolves duplicate keys by calling merge with the
// existing value and the new value, storing merge's result. An empty input
// yields an empty, non-nil map.
func ToMapMerge[E any, K comparable, V any](f func(E) (K, V), merge func(V, V) V) Collector[E, map[K]V] {
	return func(s iter.Seq[E]) map[K]V {
		out := make(map[K]V)
		for e := range s {
			k, v := f(e)
			if old, ok := out[k]; ok {
				out[k] = merge(old, v)
			} else {
				out[k] = v
			}
		}
		return out
	}
}

// Joining returns a collector that concatenates string elements with sep as
// separator. The separator is emitted only between elements, never before the
// first or after the last. An empty input yields the empty string.
func Joining(sep string) Collector[string, string] {
	return func(s iter.Seq[string]) string {
		var b strings.Builder
		first := true
		for e := range s {
			if !first {
				b.WriteString(sep)
			}
			b.WriteString(e)
			first = false
		}
		return b.String()
	}
}

// GroupingBy returns a collector that partitions elements into groups keyed by
// classifier, accumulating each group into a slice in iteration order. An empty
// input yields an empty, non-nil map.
func GroupingBy[E any, K comparable](classifier func(E) K) Collector[E, map[K][]E] {
	return func(s iter.Seq[E]) map[K][]E {
		out := make(map[K][]E)
		for e := range s {
			k := classifier(e)
			out[k] = append(out[k], e)
		}
		return out
	}
}

// GroupingByDownstream is like GroupingBy but applies downstream to each group's
// elements instead of accumulating into a slice, enabling composition. An empty
// input yields an empty, non-nil map.
//
// Group then materialize each group into a slice:
//
//	GroupingByDownstream(classifier, ToSlice[E]())
//
// Count elements per group by wrapping an xiter terminal as a Collector:
//
//	counting := collector.Collector[E, int](
//	    func(s iter.Seq[E]) int { return xiter.Size(s) },
//	)
//	GroupingByDownstream(classifier, counting)
func GroupingByDownstream[E any, K comparable, R any](classifier func(E) K, downstream Collector[E, R]) Collector[E, map[K]R] {
	return func(s iter.Seq[E]) map[K]R {
		groups := make(map[K][]E)
		for e := range s {
			k := classifier(e)
			groups[k] = append(groups[k], e)
		}
		out := make(map[K]R, len(groups))
		for k, es := range groups {
			out[k] = downstream(fromSlice(es))
		}
		return out
	}
}

// PartitioningBy returns a collector that splits elements into Pass (predicate
// true) and Fail (predicate false) halves, each preserving iteration order. An
// empty input yields a zero-value Partition with both halves nil; a side that
// receives no element is also nil.
//
//	PartitioningBy(func(e int) bool { return e%2 == 0 })(seqOf(1, 2, 3, 4, 5))
//	// Partition[int]{Pass: []int{2, 4}, Fail: []int{1, 3, 5}}
func PartitioningBy[E any](pred func(E) bool) Collector[E, Partition[E]] {
	return func(s iter.Seq[E]) Partition[E] {
		var p Partition[E]
		for e := range s {
			if pred(e) {
				p.Pass = append(p.Pass, e)
			} else {
				p.Fail = append(p.Fail, e)
			}
		}
		return p
	}
}
