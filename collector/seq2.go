package collector

import "iter"

// ToMap2 returns a collector that builds a map from a key/value sequence. On
// duplicate keys the later value overwrites the earlier one (last wins); for
// custom duplicate handling use ToMap2Merge. An empty input yields an empty,
// non-nil map.
func ToMap2[K comparable, V any]() Collector2[K, V, map[K]V] {
	return func(s iter.Seq2[K, V]) map[K]V {
		out := make(map[K]V)
		for k, v := range s {
			out[k] = v
		}
		return out
	}
}

// ToMap2Merge is like ToMap2 but resolves duplicate keys by calling merge with
// the existing value and the new value, storing merge's result. An empty input
// yields an empty, non-nil map.
func ToMap2Merge[K comparable, V any](merge func(V, V) V) Collector2[K, V, map[K]V] {
	return func(s iter.Seq2[K, V]) map[K]V {
		out := make(map[K]V)
		for k, v := range s {
			if old, ok := out[k]; ok {
				out[k] = merge(old, v)
			} else {
				out[k] = v
			}
		}
		return out
	}
}

// ToKeys returns a collector that collects keys into a slice in iteration order.
// An empty input yields a nil slice.
func ToKeys[K, V any]() Collector2[K, V, []K] {
	return func(s iter.Seq2[K, V]) []K {
		var out []K
		for k := range s {
			out = append(out, k)
		}
		return out
	}
}

// ToValues returns a collector that collects values into a slice in iteration
// order. An empty input yields a nil slice.
func ToValues[K, V any]() Collector2[K, V, []V] {
	return func(s iter.Seq2[K, V]) []V {
		var out []V
		for _, v := range s {
			out = append(out, v)
		}
		return out
	}
}
