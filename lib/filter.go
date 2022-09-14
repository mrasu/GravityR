package lib

func Filter[T any](vals []T, fn func(T) bool) []T {
	var res []T

	for _, v := range vals {
		if !fn(v) {
			res = append(res, v)
		}
	}

	return res
}
