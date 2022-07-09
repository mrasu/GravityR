package lib

func SelectF[T any](vals []T, f func(T) bool) []T {
	var res []T
	for _, v := range vals {
		if f(v) {
			res = append(res, v)
		}
	}

	return res
}
