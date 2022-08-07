package lib

func SliceOrEmpty[T any](vals []T) []T {
	if vals == nil {
		return []T{}
	}
	return vals
}
