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
