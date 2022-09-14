package lib

import (
	"golang.org/x/exp/constraints"
	"sort"
)

func SortF[T any, U constraints.Ordered](vals []T, f func(v T) U) {
	sort.Slice(vals, func(i, j int) bool {
		vi := f(vals[i])
		vj := f(vals[j])

		return vi < vj
	})
}

func SortedKeys[K constraints.Ordered, V any](vals map[K]V) []K {
	var keys []K
	for k, _ := range vals {
		keys = append(keys, k)
	}
	SortF(keys, func(k K) K {
		return k
	})

	return keys
}
