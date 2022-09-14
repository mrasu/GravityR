package lib

type Equalable[T any] interface {
	Equals(other T) bool
}

func BruteUniq[T Equalable[T]](vals []T) []T {
	var uniqVals []T
	for _, v := range vals {
		isUniq := true
		for _, uVal := range uniqVals {
			if v.Equals(uVal) {
				isUniq = false
				break
			}
		}
		if isUniq {
			uniqVals = append(uniqVals, v)
		}
	}

	return uniqVals
}

func UniqBy[T any, U comparable](vals []T, fn func(T) U) []T {
	var res []T
	found := NewSet[U]()

	for _, v := range vals {
		u := fn(v)
		if found.Contains(u) {
			continue
		}
		res = append(res, v)
		found.Add(u)
	}

	return res
}
