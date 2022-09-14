package lib

func Map[T any, U any](vals []T, fn func(T) U) []U {
	var res []U

	for _, v := range vals {
		res = append(res, fn(v))
	}

	return res
}
