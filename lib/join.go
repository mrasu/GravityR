package lib

import "strings"

func JoinF[T any](vals []T, sep string, f func(T) string) string {
	var res []string
	for _, val := range vals {
		res = append(res, f(val))
	}
	return strings.Join(res, sep)
}
